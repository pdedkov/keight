package processor_test

import (
	"context"
	"io"
	"testing"

	"keight/internal/domain"
	"keight/internal/logging"
	"keight/internal/processor"
)

type mockStorage struct {
}

func (m mockStorage) Writer(ctx context.Context, file *domain.File) (io.Writer, error) {
	return io.Discard, nil
}

func (m mockStorage) Read(ctx context.Context, id string) (*domain.File, io.Reader, error) {
	var r io.Reader
	return &domain.File{}, r, nil
}

func (m mockStorage) Remove(ctx context.Context, id string) error {
	return nil
}

func TestNew(t *testing.T) {
	_ = processor.New(&mockStorage{}, logging.Default())
}
