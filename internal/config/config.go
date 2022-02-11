package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Log   LoggerConf `mapstructure:"LOGGER"`
	DB    DBConf     `mapstructure:"DB"`
	API   APIConf    `mapstructure:"API"`
	Queue QueueConf  `mapstructure:"QUEUE"`
}

type LoggerConf struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

type DBConf struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     int    `mapstructure:"POSTGRES_PORT"`
	Name     string `mapstructure:"POSTGRES_DB"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
}

type QueueConf struct {
	Host     string `mapstructure:"RABBIT_HOST"`
	Port     int    `mapstructure:"RABBIT_PORT"`
	User     string `mapstructure:"RABBITMQ_DEFAULT_USER"`
	Password string `mapstructure:"RABBITMQ_DEFAULT_PASS"`
}

func (db *DBConf) ConnString() string {
	return fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		db.Host, strconv.Itoa(db.Port), db.User, db.Password, db.Name)
}

type APIConf struct {
	Host string `mapstructure:"API_HOST"`
	Port int    `mapstructure:"API_PORT"`
}

func (cfg *Config) bindEnv() {
	viper.AutomaticEnv()
	// TODO: this looks strange - need to understand how to unmarshal without listing ENVVARS one-by-one
	envs := []string{
		"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD",
		"STATS_INTERVAL", "API_HOST", "API_PORT",
		"RABBIT_HOST", "RABBIT_PORT", "RABBITMQ_DEFAULT_USER", "RABBITMQ_DEFAULT_PASS",
	}
	for _, key := range envs {
		_ = viper.BindEnv(key)
	}

	_ = viper.Unmarshal(&cfg.Log)
	_ = viper.Unmarshal(&cfg.DB)
	_ = viper.Unmarshal(&cfg.API)
	_ = viper.Unmarshal(&cfg.Queue)
}

func (cfg *Config) bindFile(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			_ = viper.Unmarshal(&cfg)
		}
	}
}

func New(cfgFile string) Config {
	cfg := Config{
		Log:   LoggerConf{Level: "debug"},
		DB:    DBConf{Port: 5432},
		API:   APIConf{Port: 8080},
		Queue: QueueConf{Port: 5672},
	}

	cfg.bindEnv()
	cfg.bindFile(cfgFile)

	return cfg
}
