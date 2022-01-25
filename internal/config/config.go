package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Log   LoggerConf `toml:"logger"`
	DB    DBConf     `toml:"db"`
	API   APIConf    `toml:"api"`
	Stats StatsConf  `toml:"stats"`
}

type LoggerConf struct {
	Level string `mapstructure:"LOG_LEVEL" toml:"level"`
}

type DBConf struct {
	Host     string `mapstructure:"POSTGRES_HOST" toml:"host"`
	Name     string `mapstructure:"POSTGRES_DB" toml:"name"`
	User     string `mapstructure:"POSTGRES_USER" toml:"user"`
	Password string `mapstructure:"POSTGRES_PASSWORD" toml:"password"`
}

type StatsConf struct {
	Interval int `mapstructure:"STATS_INTERVAL" toml:"interval"` // in seconds
}

func (db *DBConf) ConnString() string {
	return fmt.Sprintf(`host=%s user=%s password=%s dbname=%s sslmode=disable`,
		db.Host, db.User, db.Password, db.Name)
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
