package dynamic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	ErrTunnelNotExits = errors.New("tunnel not exits")
	ErrDownloadFailed = errors.New("download failed")
)

func isTunnelNotExist(err error) bool {
	return errors.Is(err, ErrTunnelNotExits)
}

type Remote interface {
	Sync(name string) error
}

func NewRemote() Remote {
	if remote == "" {
		return nil
	}

	u, err := url.Parse(remote)
	if err != nil {
		log.Panicf("parsing remote url error: %v", err)
	}

	switch u.Scheme {
	case "s3":
		return NewS3Remote(u.Host)
	default:
		log.Panicf("unknown remote scheme: %s", u.Scheme)
	}

	return nil
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
		return nil, fmt.Errorf("failed to create client, %w", err)
	}

	return s3.NewFromConfig(cfg), nil
}

func (r *S3Remote) downloadFileFromS3(remoteFilePath string, localFilePath string) error {
	client, err := r.createS3Client()
	if err != nil {
		return fmt.Errorf("failed to create s3 client, %w", err)
	}

	// Create a file to write the S3 Object contents to.
	getObjectResponse, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(remoteFilePath),
	})

	if err != nil {
		log.Printf("failed to get object, %v", err)
		return ErrTunnelNotExits
	}

	file, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %w", localFilePath, err)
	}
	defer file.Close()

	written, err := io.Copy(file, getObjectResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to write file contents! %w", err)
	} else if written != *getObjectResponse.ContentLength {
		return fmt.Errorf("wrote a different size than was given to us")
	}

	return nil
}

func (r *S3Remote) batchDownloadFilesFromS3(name string) error {
	files := []string{
		fmt.Sprintf("libcgo_%s.so", name),
		fmt.Sprintf("libgo_%s.so", name),
	}

	var wg sync.WaitGroup
	var errChan = make(chan error, len(files))
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			localFilePath := filepath.Join(filepath.Join(GetWarehouse(), runtime.Version(), name), file)
			remoteFilePath := filepath.ToSlash(filepath.Join(runtime.Version(), name, file))

			if stat, err := os.Stat(localFilePath); err != nil {
				if os.IsNotExist(err) {
					log.Printf("%s not found, downloading from s3://%s...", localFilePath, filepath.Join(r.bucket, remoteFilePath))
					if err := r.downloadFileFromS3(remoteFilePath, localFilePath); err != nil {
						log.Printf("failed to download file from s3, %v", err)
						errChan <- err
						return
					}
				} else {
					log.Printf("failed to stat file, %v", err)
					errChan <- err
					return
				}
			} else if stat.Size() == 0 {
				log.Printf("%s is empty, downloading from s3://%s...", localFilePath, filepath.Join(r.bucket, remoteFilePath))
				if err := os.Remove(localFilePath); err != nil {
					log.Printf("failed to remove file, %v", err)
					errChan <- err
					return
				}
				if err := r.downloadFileFromS3(remoteFilePath, localFilePath); err != nil {
					log.Printf("failed to download file from s3, %v", err)
					errChan <- err
					return
				}
			} else {
				log.Printf("%s is already exists", localFilePath)
			}
		}(file)
	}
	wg.Wait()

	if len(errChan) > 0 {
		log.Printf("%d errors occurred during downloading", len(errChan))
		for err := range errChan {
			if isTunnelNotExist(err) {
				return ErrTunnelNotExits
			}
		}
		return ErrDownloadFailed
	}

	return nil
}

func (r *S3Remote) Sync(name string) error {
	dir := filepath.Join(GetWarehouse(), runtime.Version(), name)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create dir %s, %w", dir, err)
			}
		} else {
			return fmt.Errorf("failed to stat dir, %w", err)
		}
	}

	startTime := time.Now()
	if err := r.batchDownloadFilesFromS3(name); err != nil {
		os.RemoveAll(dir)
		if isTunnelNotExist(err) {
			return ErrTunnelNotExits
		}
		return fmt.Errorf("failed to download files from s3, %w", err)
	}
	log.Printf("download files from s3 took %v", time.Since(startTime))

	return nil
}
