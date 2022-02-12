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

type Tasks struct {
	Addr        string `mapstructure:"tasks_addr" default:"localhost:6379"`
	Concurrency int    `mapstructure:"tasks_concurrency" validate:"min=1" default:"50"`
}

type App struct {
	Secret   string `mapstructure:"app_secret"`
	SeedPath string `mapstructure:"app_seed_path"`
}

type Log struct {
	Level string `mapstructure:"log_level"`
}

type RedisBackend struct {
	Addr string `mapstructure:"feedback_backend__redis__addr"`
}

type FeedbackBackeds struct {
	Redis RedisBackend `mapstructure:",squash"`
}

type Config struct {
	DB               DB              `mapstructure:",squash"`
	Tasks            Tasks           `mapstructure:",squash"`
	Log              Log             `mapstructure:",squash"`
	App              App             `mapstructure:",squash"`
	FeedbackBackends FeedbackBackeds `mapstructure:",squash"`
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
