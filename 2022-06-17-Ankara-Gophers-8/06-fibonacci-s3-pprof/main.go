package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var profilers = []string{
	"heap",
	"allocs",
	"goroutine",
	"threadcreate",
	"block",
	"mutex",
	"cpu",
}

var logger = log.New(os.Stdout, "", log.LstdFlags)

func fibonacciWithMemory(n uint64) uint64 {
	cache := make([]uint64, n+1)
	cache[0] = 0
	cache[1] = 1
	var i uint64
	for i = 2; i <= n; i++ {
		cache[i] = cache[i-1] + cache[i-2]
	}

	return cache[n]
}

func fibonacciWithRecursion(n uint64) uint64 {
	if n <= 1 {
		return n
	}
	return fibonacciWithRecursion(n-1) + fibonacciWithRecursion(n-2)
}

func getCreds() (bucket, region, accessKey, secretKey string, ok bool) {
	ok = true
	bucket = os.Getenv("ANKGOPHERS_DEMO_S3_BUCKET")
	region = os.Getenv("ANKGOPHERS_DEMO_AWS_REGION")
	accessKey = os.Getenv("ANKGOPHERS_DEMO_AWS_ACCESS_KEY_ID")
	secretKey = os.Getenv("ANKGOPHERS_DEMO_AWS_SECRET_ACCESS_KEY")
	if accessKey == "" || secretKey == "" || region == "" || bucket == "" {
		ok = false
	}
	return
}

func newS3Client(config *aws.Config) (*s3.S3, error) {
	session, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("could not initialize new aws session: %v", err)
	}
	return s3.New(session), nil
}

func uploadFileToS3(cl *s3.S3, bucket, key string, data []byte) error {
	if bucket == "" {
		return fmt.Errorf("bucket is required")
	}
	b := bytes.NewReader(data)
	_, err := cl.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		Body:                 b,
		ContentLength:        aws.Int64(int64(b.Len())),
		ContentType:          aws.String("application/octet-stream"),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func captureProfile(ctx context.Context, name string) (*bytes.Buffer, error) {
	if name == "cpu" {
		if runtime.GOOS == "windows" {
			return nil, fmt.Errorf("CPU profiling is not supported on windows")
		}
		b := new(bytes.Buffer)
		if err := pprof.StartCPUProfile(b); err != nil {
			return nil, fmt.Errorf("failed to start CPU profile: %v", err)
		}
		select {
		case <-time.After(10 * time.Second):
		case <-ctx.Done():
		}
		pprof.StopCPUProfile()
		return b, nil
	}
	p := pprof.Lookup(name)
	if p == nil {
		return nil, fmt.Errorf("unknown profile: %v", name)
	}
	b := new(bytes.Buffer)
	if err := p.WriteTo(b, 0); err != nil {
		return nil, fmt.Errorf("could not save dump to buffer: %v", err)
	}
	return b, nil
}

func startContinuousProfiling(ctx context.Context, interval time.Duration, bucket string, awsConfig *aws.Config) error {
	cl, err := newS3Client(awsConfig)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(interval)
	go func() {
		logger.Printf("Starting profiling with %v\n", profilers)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				logger.Println("Stopped profiler")
				return
			case <-ticker.C:
				for _, p := range profilers {
					logger.Printf("Capturing %s profile\n", p)
					b, err := captureProfile(ctx, p)
					if err != nil {
						logger.Fatalf("Failed to capture %s profile, err: %v\n", p, err)
					}
					now := time.Now().UTC()
					key := fmt.Sprintf("%s/%s_%s.pb.gz", now.Format("2006-01-02"), now.Format("15-04-05"), p)
					if err := uploadFileToS3(cl, bucket, key, b.Bytes()); err != nil {
						logger.Fatalf("Failed to upload %s profile, err: %v\n", p, err)
					}
					logger.Printf("Committed %s profile to %s\n", p, key)
				}
			}
		}
	}()
	return nil
}

var (
	interval = flag.Duration("interval", 10*time.Second, "interval between profiling dumps")
)

func main() {
	bucket, region, accessKey, secretKey, ok := getCreds()
	if !ok {
		logger.Fatalf("AWS credentials are not passed properly")
	}
	startContinuousProfiling(context.Background(), *interval, bucket, &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	go func() {
		for {
			_ = fibonacciWithMemory(10000000)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	go func() {
		for {
			_ = fibonacciWithRecursion(45)
		}
	}()
	select {}
}
