package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"keight/internal/domain"
	"keight/internal/helper"
	"keight/internal/logging"

	"github.com/google/uuid"
)

type Proccessor interface {
	Upload(ctx context.Context, r *multipart.Reader, size int64) (*domain.File, error)
	Download(ctx context.Context, id string) (*domain.File, io.Reader, error)
}

type Handler interface {
	Handle(pattern string, handler http.Handler)
}

// API is http api server handlers
type API struct {
	log  logging.Logger
	proc Proccessor
}

// New creates new api service
func New(proc Proccessor, log logging.Logger) *API {
	return &API{
		proc: proc,
		log:  log,
	}
}

func (a *API) ApplyHanders(router Handler) {
	router.Handle("POST /upload", helper.Recoverer(a.log)(http.HandlerFunc(a.upload)))
	router.Handle("GET /download", helper.Recoverer(a.log)(http.HandlerFunc(a.download)))
}

func (a *API) upload(w http.ResponseWriter, r *http.Request) {
	log := a.log.WithContext(r.Context())
	log.Info("process upload")

	reader, err := r.MultipartReader()
	if err != nil {
		log.WithError(err).Error("cant get multipart form reader from request")
		helper.WrapError(w, err, http.StatusBadRequest)
		return
	}
	file, err := a.proc.Upload(logging.NewContext(r.Context(), log), reader, r.ContentLength)
	if err != nil {
		log.WithError(err).Error("upload multipart failed")
		helper.WrapError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(file); err != nil {
		log.WithError(err).WithField(logging.ID, file.ID).Error("cant upload file")
		return
	}
	log.WithField(logging.ID, file.ID).WithField(logging.Size, file.Size).Info("upload ok")
}

func (a *API) download(w http.ResponseWriter, r *http.Request) {
	log := a.log.WithContext(r.Context())
	id := r.URL.Query().Get("id")
	log = log.WithField(logging.ID, id)
	if _, err := uuid.Parse(id); err != nil {
		log.WithError(err).Error("wrong file id")
		helper.WrapError(w, err, http.StatusBadRequest)
		return
	}

	file, reader, err := a.proc.Download(logging.NewContext(r.Context(), log), id)
	if err != nil {
		log.WithError(err).Error("cant download file")
		helper.WrapError(w, err, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename=%q`, file.Name))
	w.WriteHeader(http.StatusOK)
	ln, err := io.Copy(w, reader)
	log = log.WithField(logging.Size, ln)
	if err != nil {
		log.WithError(err).Error("cant download file")
		return
	}
	log.Info("downloaded ok")
}
