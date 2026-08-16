package services

import (
	"log"
	"os"
	"strings"
)

type FileStore interface {
	DownloadFile(bucketName, key, fileName string) error
	UploadFile(bucketName, key, fileName string) error
}

type NoopStore struct{}

func NewNoopStore() NoopStore {
	return NoopStore{}
}

func (NoopStore) DownloadFile(_, _, _ string) error {
	return nil
}

func (NoopStore) UploadFile(_, _, _ string) error {
	return nil
}

func awsDisabled() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("AWS_DISABLED")))
	return v == "1" || v == "true" || v == "yes"
}

func sqlitePath() string {
	if path := os.Getenv("SQLITE_PATH"); path != "" {
		return path
	}
	return "/tmp/todos.db"
}

func ConfigStorage() (FileStore, *S3Details) {
	details := &S3Details{
		Key:      "tmp/todos.db",
		FileName: sqlitePath(),
	}
	if awsDisabled() || os.Getenv("BUCKET_NAME") == "" {
		log.Println("AWS disabled; skipping S3")
		return NewNoopStore(), details
	}

	s3Svc, awsDetails := ConfigAws()
	awsDetails.FileName = details.FileName
	return s3Svc, awsDetails
}
