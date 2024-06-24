package db

import (
	"context"
	"encoding/json"
	"fmt"

	"keight/internal/domain"
	"keight/internal/logging"

	"github.com/go-redis/redis/v8"
)

type configuration interface {
	Address() string
	DB() int
	Password() string
}

const (
	filePrefix = "file"
)

type Redis struct {
	log    logging.Logger
	client *redis.Client
}

// NewRedis creates new redis storage instance
func NewRedis(cfg configuration, log logging.Logger) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		DB:       cfg.DB(),
		Password: cfg.Password(),
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	log.WithField(logging.Address, cfg.Address()).WithField(logging.DB, cfg.DB()).
		Info("connected to redis instance")
	return &Redis{
		log:    log,
		client: client,
	}, nil
}

// InsertFile saves new file info
func (r *Redis) InsertFile(ctx context.Context, file *domain.File) error {
	val, err := json.Marshal(file)
	if err != nil {
		r.log.WithError(err).Error("cant encode file")
		return err
	}
	if err = r.client.Set(ctx, r.fileKey(file.ID), val, 0).Err(); err != nil {
		r.log.WithError(err).Error("cant insert file in db")
		return err
	}
	r.log.WithField(logging.ID, file.ID).WithField(logging.Filename, file.Name).
		Info("file inserted")
	return nil
}

func (r *Redis) fileKey(id string) string {
	return fmt.Sprintf("%s:%s", filePrefix, id)
}

// GetFile returns file by its id
func (r *Redis) GetFile(ctx context.Context, id string) (*domain.File, error) {
	log := r.log.WithField(logging.ID, id)
	resp := r.client.Get(ctx, r.fileKey(id))
	if resp.Err() != nil {
		err := resp.Err()
		log.WithError(err).Error("cant get file from db")
		return nil, err
	}
	raw, err := resp.Bytes()
	if err != nil {
		log.WithError(err).Error("cant load from db")
		return nil, err
	}
	var file domain.File
	err = json.Unmarshal(raw, &file)
	if err != nil {
		log.WithError(err).Error("cant encode data to file")
		return nil, err
	}
	file.ID = id
	return &file, nil
}

// RemoveFile removes file from DB
func (r *Redis) RemoveFile(ctx context.Context, id string) error {
	log := r.log.WithField(logging.ID, id)
	err := r.client.Del(ctx, r.fileKey(id)).Err()
	if err != nil {
		log.WithError(err).Error("cant remove file from db")
		return err
	}
	return nil
}
