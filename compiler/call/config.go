package call

type Config struct {
	BasePath string
}

func (cfg *Config) SetBasePath(path string) {
	cfg.BasePath = path
}
