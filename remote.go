package dynamic

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Remote interface {
	ExistsOrSync(name string)
}

var (
	remote Remote
)

func init() {
	if s, ok := os.LookupEnv("DYNAMIC_REMOTE"); !ok {
		return
	} else {
		u, err := url.Parse(s)
		if err != nil {
			log.Fatalf("parsing remote url error: %v", err)
		}
		switch u.Scheme {
		case "s3":
			remote = NewS3Remote(u.Host)
		default:
			log.Fatalf("unknown remote scheme: %s", u.Scheme)
		}
	}
}

type S3Remote struct {
	bucket string
}

func NewS3Remote(bucket string) *S3Remote {
	return &S3Remote{
		bucket: bucket,
	}
}

func (r *S3Remote) createS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create client, %v", err)
	}

	return s3.NewFromConfig(cfg), nil
}

func (r *S3Remote) downloadFileFromS3(remoteFilePath string, localFilePath string) error {
	client, err := r.createS3Client()
	if err != nil {
		log.Fatalf("failed to create s3 client, %v", err)
	}

	// Create a file to write the S3 Object contents to.
	getObjectResponse, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(remoteFilePath),
	})

	if err != nil {
		return fmt.Errorf("failed to get object, %v", err)
	}

	file, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", localFilePath, err)
	}
	defer file.Close()

	written, err := io.Copy(file, getObjectResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to write file contents! %v", err)
	} else if written != getObjectResponse.ContentLength {
		return fmt.Errorf("wrote a different size than was given to us")
	}

	return nil
}

// TODO: passing errors out, should not use fatal
func (r *S3Remote) batchDownloadFilesFromS3(name string) {
	files := []string{
		fmt.Sprintf("libcgo_%s.so", name),
		fmt.Sprintf("libgo_%s.so", name),
	}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			localFilePath := filepath.Join(filepath.Join(warehouse, name), file)
			remoteFilePath := filepath.ToSlash(filepath.Join(name, file))

			if stat, err := os.Stat(localFilePath); err != nil {
				if os.IsNotExist(err) {
					log.Printf("%s does not exist, downloading from s3[%s]...", localFilePath, remoteFilePath)
					if err := r.downloadFileFromS3(remoteFilePath, localFilePath); err != nil {
						log.Fatalf("failed to download file from s3, %v", err)
					}
				} else {
					log.Fatalf("failed to stat file, %v", err)
				}
			} else if stat.Size() == 0 {
				log.Printf("%s is empty, downloading from s3[%s]...", localFilePath, remoteFilePath)
				if err := os.Remove(localFilePath); err != nil {
					log.Fatalf("failed to remove file, %v", err)
				}
				if err := r.downloadFileFromS3(remoteFilePath, localFilePath); err != nil {
					log.Fatalf("failed to download file from s3, %v", err)
				}
			} else {
				log.Printf("%s is already exists", localFilePath)
			}
		}(file)
	}
	wg.Wait()
}

func (r *S3Remote) ExistsOrSync(name string) {
	dir := filepath.Join(warehouse, name)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Fatalf("failed to create dir %s, %v", dir, err)
			}
		} else {
			log.Fatalf("failed to stat dir, %v", err)
		}
	}

	startTime := time.Now()
	r.batchDownloadFilesFromS3(name)
	log.Printf("download files from s3 took %v", time.Since(startTime))
}
