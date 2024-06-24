package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Mime     string    `json:"mime"`
	Storages Storage   `json:"storages"`
	Created  time.Time `json:"created"`
}

func NewFile(name, mime string, size int64) *File {
	return &File{
		ID:      uuid.New().String(),
		Name:    name,
		Mime:    mime,
		Size:    size,
		Created: time.Now(),
	}
}
