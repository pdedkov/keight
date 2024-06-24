package config

type Storage struct {
	Count       int `envconfig:"STORAGE_COUNT" default:"6"`
	ChunksCount int `envconfig:"STORAGE_CHUNKS" default:"6"`
}
