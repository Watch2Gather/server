package config

type (
	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}
	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port int    `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}
	WS struct {
		WSHost string `env-required:"true" yaml:"ws_host" env:"WS_PORT"`
		WSPort int    `env-required:"true" yaml:"ws_port" env:"WS_PORT"`
	}
	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}
)
