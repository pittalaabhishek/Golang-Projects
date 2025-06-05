package storage

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	bolt "go.etcd.io/bbolt"
)

const (
	bucketName = "tasks"
	dbFileName = "tasks.db"
	appDir 	   = ".task"
)

type TaskStore struct {
	db *bolt.DB
}

func NewTaskStore() (*TaskStore, error) {
	dbPath, err := getDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &TaskStore{db: db}

	// Initialize bucket
	if err := store.initBucket(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize bucket: %w", err)
	}

	return store, nil
}

func (s *TaskStore) initBucket() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
}

func (s *TaskStore) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func getDatabasePath() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, appDir, dbFileName), nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (s *TaskStore) AddTask(description string) (*Task, error) {
	var task *Task

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}

		id, err := bucket.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to generate ID: %w", err)
		}

		task = &Task{
			ID:          int(id),
			Description: description,
			Completed:   false,
			CreatedAt:   time.Now(),
		}

		data, err := task.MarshalBinary()
		if err != nil {
			return fmt.Errorf("failed to marshal task: %w", err)
		}

		key := itob(task.ID)
		return bucket.Put(key, data)
	})

	return task, err
}

func (s *TaskStore) GetIncompleteTasks() ([]Task, error) {
	var tasks []Task

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return nil // No tasks yet
		}

		return bucket.ForEach(func(k, v []byte) error {
			var task Task
			if err := task.UnmarshalBinary(v); err != nil {
				return fmt.Errorf("failed to unmarshal task: %w", err)
			}

			if !task.Completed {
				tasks = append(tasks, task)
			}

			return nil
		})
	})

	return tasks, err
}

func (s *TaskStore) CompleteTask(position int) (*Task, error) {
	incompleteTasks, err := s.GetIncompleteTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to get incomplete tasks: %w", err)
	}

	if position < 1 || position > len(incompleteTasks) {
		return nil, fmt.Errorf("invalid task number: %d (must be between 1 and %d)", position, len(incompleteTasks))
	}

	task := incompleteTasks[position-1]
	task.Completed = true
	now := time.Now()
	task.CompletedAt = &now

	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}

		data, err := task.MarshalBinary()
		if err != nil {
			return fmt.Errorf("failed to marshal task: %w", err)
		}

		key := itob(task.ID)
		return bucket.Put(key, data)
	})

	return &task, err
}

func (s *TaskStore) DeleteTask(position int) (*Task, error) {
	incompleteTasks, err := s.GetIncompleteTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to get incomplete tasks: %w", err)
	}

	if position < 1 || position > len(incompleteTasks) {
		return nil, fmt.Errorf("invalid task number: %d (must be between 1 and %d)", position, len(incompleteTasks))
	}

	task := incompleteTasks[position-1]

	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}

		key := itob(task.ID)
		return bucket.Delete(key)
	})

	return &task, err
}

func (s *TaskStore) GetCompletedTasksToday() ([]Task, error) {
	var tasks []Task

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var task Task
			if err := task.UnmarshalBinary(v); err != nil {
				return fmt.Errorf("failed to unmarshal task: %w", err)
			}

			if task.IsCompletedToday() {
				tasks = append(tasks, task)
			}

			return nil
		})
	})

	return tasks, err
}