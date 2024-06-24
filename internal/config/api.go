package config

type API struct {
	Port string `envconfig:"HTTP_API_PORT" required:"true"`
}

func (a *API) Addr() string {
	return ":" + a.Port
}
