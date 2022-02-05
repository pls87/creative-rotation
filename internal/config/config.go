package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Log   LoggerConf `mapstructure:"LOGGER"`
	DB    DBConf     `mapstructure:"DB"`
	API   APIConf    `mapstructure:"API"`
	Stats StatsConf  `mapstructure:"STATS"`
}

type LoggerConf struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

type DBConf struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	Name     string `mapstructure:"POSTGRES_DB"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
}

type StatsConf struct {
	Interval int `mapstructure:"STATS_INTERVAL"` // in seconds
}

func (db *DBConf) ConnString() string {
	return fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		db.Host, db.Port, db.User, db.Password, db.Name)
}

type APIConf struct {
	Host string `mapstructure:"API_HOST" toml:"host"`
	Port int    `mapstructure:"API_PORT" toml:"port"`
}

func (cfg *Config) bindEnv() {
	viper.AutomaticEnv()
	// TODO: this looks strange - need to understand how to unmarshal without listing ENVVARS one-by-one
	_ = viper.BindEnv("POSTGRES_HOST")
	_ = viper.BindEnv("POSTGRES_DB")
	_ = viper.BindEnv("POSTGRES_USER")
	_ = viper.BindEnv("POSTGRES_PASSWORD")
	_ = viper.BindEnv("STATS_INTERVAL")
	_ = viper.BindEnv("API_HOST")
	_ = viper.BindEnv("API_PORT")
	_ = viper.Unmarshal(&cfg.Log)
	_ = viper.Unmarshal(&cfg.DB)
	_ = viper.Unmarshal(&cfg.API)
	_ = viper.Unmarshal(&cfg.Stats)
}

func (cfg *Config) bindFile(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			_ = viper.Unmarshal(cfg)
		}
	}
}

func New(cfgFile string) Config {
	cfg := Config{
		Log:   LoggerConf{Level: "debug"},
		DB:    DBConf{},
		API:   APIConf{Port: 8083},
		Stats: StatsConf{Interval: 1},
	}

	cfg.bindEnv()
	cfg.bindFile(cfgFile)

	return cfg
}
