package conf

import (
	"github.com/spf13/viper"
)

type DB struct {
	Host     string `mapstructure:"db_host"`
	Port     int32  `mapstructure:"db_port"`
	Database string `mapstructure:"db_database"`
	Username string `mapstructure:"db_username"`
	Password string `mapstructure:"db_password"`
	Options  string `mapstructure:"db_options"`
}

type Worker struct {
	BackendAddr     string `mapstructure:"workerBackendAddr" default:"localhost:6379"`
	Concurrency     int    `mapstructure:"workerConcurrency" validate:"min=1" default:"50"`
	FeedBackendAddr string `mapstructure:"feedBackendAddr" default:"localhost:6379"`
}

type FeedProcessor struct {
}

type FeedBackend struct {
	Addr      string `mapstructure:"feedBackendAddr" default:"localhost:9090"`
	CacheAddr string `mapstructure:"workerBackendAddr" default:"localhost:6379"`
	DB        DB     `mapstructure:",squash"`
}

type App struct {
	Secret   string `mapstructure:"app_secret"`
	SeedPath string `mapstructure:"app_seed_path"`
}

type Log struct {
	Level string `mapstructure:"log_level"`
}

type Config struct {
	DB          DB          `mapstructure:",squash"`
	Worker      Worker      `mapstructure:",squash"`
	Log         Log         `mapstructure:",squash"`
	App         App         `mapstructure:",squash"`
	FeedBackend FeedBackend `mapstructure:",squash"`
}

func LoadConfig(args ...string) (config Config, err error) {
	envfile := ".env"
	if len(args) > 0 {
		envfile = args[0]
	}
	viper.SetConfigFile(envfile)

	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(&config)
	return config, err
}
