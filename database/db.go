package database

import (
	"davidc/todo-api/services"
	"log"
	"sync"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

const schema = `
CREATE TABLE IF NOT EXISTS todos (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task TEXT NOT NULL,
  completed INTEGER NOT NULL DEFAULT 0
);
`

var (
	s3db    *sqlx.DB
	store   services.FileStore
	details *services.S3Details
	mu      sync.RWMutex
)

func InitDb(fileStore services.FileStore, s3Details *services.S3Details) error {
	if s3db != nil {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	if s3db != nil {
		return nil
	}
	store = fileStore
	details = s3Details
	if err := fileStore.DownloadFile(s3Details.BucketName, s3Details.Key, s3Details.FileName); err != nil {
		log.Printf("DB download skipped: %v", err)
	}
	newDb, err := sqlx.Open("sqlite", s3Details.FileName)
	if err != nil {
		return err
	}
	if _, err := newDb.Exec(schema); err != nil {
		return err
	}
	s3db = newDb
	return nil
}

func Write(fn func(db *sqlx.DB) error) {
	mu.Lock()
	defer mu.Unlock()
	err := fn(s3db)
	if err != nil {
		return
	}
	if store == nil || details == nil {
		return
	}
	if err := store.UploadFile(details.BucketName, details.Key, details.FileName); err != nil {
		log.Printf("DB upload failed: %v", err)
	}
}

func Read(fn func(db *sqlx.DB)) {
	mu.RLock()
	defer mu.RUnlock()
	fn(s3db)
}
