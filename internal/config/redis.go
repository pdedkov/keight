package config

type Redis struct {
	RedisAddress  string `envconfig:"REDIS_ADDRESS" required:"true"`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
}

func (r *Redis) Address() string {
	return r.RedisAddress
}

func (r *Redis) DB() int {
	return r.RedisDB
}

func (r *Redis) Password() string {
	return r.RedisPassword
}
