package dynamic

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Config struct {
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

type Remote struct {
	*S3Config
}

func NewRemote(config *S3Config) *Remote {
	return &Remote{config}
}

var (
	remote = NewRemote(
		&S3Config{
			Region:    os.Getenv("S3_REGION"),
			Bucket:    os.Getenv("S3_BUCKET"),
			AccessKey: os.Getenv("S3_ACCESS_KEY"),
			SecretKey: os.Getenv("S3_SECRET_KEY"),
		},
	)
)

func (r *Remote) createS3Session() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(r.Region),
		Credentials: credentials.NewStaticCredentials(r.AccessKey, r.SecretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session, %v", err)
	}

	return sess, nil
}

func (r *Remote) downloadFileFromS3(remoteFilePath string, localFilePath string) error {
	sess, err := r.createS3Session()
	if err != nil {
		log.Fatalf("failed to create s3 session, %v", err)
	}

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", localFilePath, err)
	}

	s3Client := s3.New(sess)
	downloader := s3manager.NewDownloaderWithClient(s3Client, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024 // 64MB per part
	})

	// Write the contents of S3 Object to the file
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(remoteFilePath),
	})
	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}

	return nil
}

func (r *Remote) batchDownloadFilesFromS3(name string) {
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
			remoteFilePath := filepath.Join(name, file)

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

func MustExists(path string) {
	remote.MustExists(path)
}

func (r *Remote) MustExists(name string) {
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
