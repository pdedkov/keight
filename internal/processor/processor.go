package processor

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"sync"

	"keight/internal/domain"
	"keight/internal/logging"
)

const (
	contentType = "Content-Type"
	// from multipart lib
	maxMIMEHeaderSize = 10 << 20
)

type Storager interface {
	Writer(ctx context.Context, file *domain.File) (io.Writer, error)
	Read(ctx context.Context, id string) (*domain.File, io.Reader, error)
	Remove(ctx context.Context, id string) error
}

type Processor struct {
	log     logging.Logger
	storage Storager
	mu      sync.RWMutex
}

// New creates new multipart processing service
func New(storage Storager, log logging.Logger) *Processor {
	return &Processor{
		log:     log,
		storage: storage,
	}
}

// Upload processing multipart file uploading, processing file and upload it to storage cluster
func (p *Processor) Upload(ctx context.Context, r *multipart.Reader, size int64) (*domain.File, error) {
	var (
		log         = logging.FromContext(ctx)
		buf         = new(bytes.Buffer)
		fileName    string
		bytesCopied int64
		file        *domain.File
	)

	// because of dummy storage implementation just lock/rlock upload/download operations
	p.mu.Lock()
	defer p.mu.Unlock()

	var processed bool
	for {
		part, err := r.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Info("processing eof")
				break
			}
			log.WithError(err).Error("cant get multipart next part from request")
			return nil, err
		}
		// process only first passed file but read all request parts
		if processed {
			continue
		}
		fileName = part.FileName()
		if fileName == "" {
			// trying to calculate real file size by content-lenght - all parts without files size
			partSize, _ := io.Copy(io.Discard, part)
			size -= partSize
			log.WithField(logging.Formname, part.FormName()).Debug("skip part")
			continue
		}
		file = domain.NewFile(fileName, part.Header.Get(contentType), size)
		// trying to remove part headers
		w, err := p.storage.Writer(ctx, file)
		if err != nil {
			log.WithError(err).Error("cant process storage")
			return nil, err
		}
		bytesCopied, err = io.Copy(w, part)
		if err != nil {
			log.WithError(err).Error("cant copy from part to storage")
			return nil, err
		}
		log = log.WithField(logging.Formname, fileName)
		if _, err := io.Copy(buf, part); err != nil {
			log.WithError(err).Error("cant copy part")
			return nil, err
		}
		// process only first passed file
		processed = true
	}
	// aprox value for simplicity
	// TODO real file size without reading it's content (remove all headers?)
	if file != nil && bytesCopied+maxMIMEHeaderSize < size {
		err := domain.ErrFileSizeMismatch
		log.WithError(err).WithField(logging.Size, size).
			WithField(logging.Len, bytesCopied).
			WithField(logging.Formname, fileName).Warn("file size mismatch")
		if errRemove := p.storage.Remove(ctx, file.ID); errRemove != nil {
			log.WithError(errRemove).Warn("cant remove file")
		}
		return nil, err
	}
	log.WithField(logging.Size, bytesCopied).Info("uploading done")
	return file, nil
}

// Download downloads file from by file ident
func (p *Processor) Download(ctx context.Context, id string) (*domain.File, io.Reader, error) {
	log := logging.FromContext(ctx)

	p.mu.RLock()
	defer p.mu.RUnlock()

	file, r, err := p.storage.Read(ctx, id)
	if err != nil {
		log.WithError(err).Error("cant read file from storage")
		return nil, nil, err
	}
	log.Info("downloading done")
	return file, r, nil
}
