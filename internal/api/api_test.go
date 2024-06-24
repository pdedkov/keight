package api_test

import (
	"context"
	"io"
	"mime/multipart"
	"testing"

	"keight/internal/api"
	"keight/internal/domain"
	"keight/internal/logging"
)

type mockProcessor struct {
}

func (m mockProcessor) Upload(ctx context.Context, r *multipart.Reader, size int64) (*domain.File, error) {
	return &domain.File{}, nil
}

func (m mockProcessor) Download(ctx context.Context, id string) (*domain.File, io.Reader, error) {
	var r io.Reader
	return &domain.File{}, r, nil
}

func TestNew(t *testing.T) {
	_ = api.New(&mockProcessor{}, logging.Default())
}
