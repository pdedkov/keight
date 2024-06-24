//go:build integration

package db_test

import (
	"context"
	"testing"

	"keight/internal/config"
	"keight/internal/db"
	"keight/internal/domain"
	"keight/internal/logging"

	"github.com/kelseyhightower/envconfig"
)

const (
	testFileID = "bb0e138e-6ac2-4e4e-a755-8ea05bbf5fbd"
)

func TestNewRedis(t *testing.T) {
	_ = newRedisIntegrationClient(t)
}

func TestRedis_InsertFile(t *testing.T) {
	drv := newRedisIntegrationClient(t)
	file := domain.NewFile("tests.dmg", "text/plain", 100)
	file.Storages = []int{1, 2, 3, 4, 5}
	err := drv.InsertFile(context.Background(), file)
	if err != nil {
		t.Error(err)
	}
	t.Log("ok")
}

func TestRedis_GetFile(t *testing.T) {
	drv := newRedisIntegrationClient(t)
	file, err := drv.GetFile(context.Background(), testFileID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", file)
}

func newRedisIntegrationClient(t *testing.T) *db.Redis {
	t.Helper()
	cfg := &config.Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		t.Fatal(err)
	}
	log := logging.NewSLog("", cfg.Log)
	store, err := db.NewRedis(cfg.Redis, log)
	if err != nil {
		t.Fatal(err)
	}
	return store
}
