package config

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  struct {
		Address string `yaml:"address"`
	} `yaml:"http_server"`
}
