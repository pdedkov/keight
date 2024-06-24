package storage

import (
	"context"
	"io"

	"keight/internal/domain"
	"keight/internal/logging"
)

type Databaser interface {
	InsertFile(ctx context.Context, file *domain.File) error
	GetFile(ctx context.Context, id string) (*domain.File, error)
	RemoveFile(ctx context.Context, id string) error
}
type Storage struct {
	log       logging.Logger
	instances []map[string][]byte
	// number of storage chunks (hosts)
	count int
	// uploaded file chunk count
	chunksCount int
	db          Databaser
	// index of next start chunk
	next int
}

func New(count, chunksCount int, db Databaser, log logging.Logger) *Storage {
	return &Storage{
		log:         log,
		db:          db,
		instances:   make([]map[string][]byte, count),
		count:       count,
		chunksCount: chunksCount,
	}
}

func (s *Storage) Writer(ctx context.Context, file *domain.File) (io.Writer, error) {
	hosts := make([]map[string][]byte, s.count)
	// select N servers for file chunks round-robin + N next hosts
	for storageIndex, i := s.next, 0; i < s.chunksCount; i, storageIndex = i+1, storageIndex+1 {
		if storageIndex == s.count {
			storageIndex = 0
		}
		file.Storages = append(file.Storages, storageIndex)
		if s.instances[storageIndex] == nil {
			s.instances[storageIndex] = make(map[string][]byte)
		}
		hosts[i] = s.instances[storageIndex]
	}
	// increment next writing chunk index
	s.next++
	// s.next locked by processor layout
	if s.next >= s.count {
		s.next = 0
	}
	// save file to db
	if err := s.db.InsertFile(ctx, file); err != nil {
		s.log.WithError(err).Error("cant insert file")
		return nil, err
	}
	// write file info in db and prepare writer
	return newWriter(file.ID, file.Size/int64(s.chunksCount), hosts...), nil
}

func (s *Storage) Read(ctx context.Context, id string) (*domain.File, io.Reader, error) {
	file, err := s.db.GetFile(ctx, id)
	if err != nil {
		s.log.WithError(err).Error("cant get file for reading")
		return nil, nil, err
	}

	hosts := make([]map[string][]byte, len(file.Storages))
	for i := range file.Storages {
		hosts[i] = s.instances[file.Storages[i]]
	}
	return file, newReader(id, hosts...), nil
}

// Remove removes file from storage and database
func (s *Storage) Remove(ctx context.Context, id string) error {
	file, _, err := s.Read(ctx, id)
	if err != nil {
		s.log.WithError(err).Error("cant get file for remove")
		return err
	}
	// remove chunks from storages
	for _, storageIndex := range file.Storages {
		delete(s.instances[storageIndex], id)
	}
	// remove file data from db
	return s.db.RemoveFile(ctx, file.ID)
}
