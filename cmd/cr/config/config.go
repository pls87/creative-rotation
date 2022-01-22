package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf `toml:"logger"`
	DB     DBConf     `toml:"db"`
	API    APIConf    `toml:"api"`
}

type LoggerConf struct {
	Level string `mapstructure:"LOGGERLEVEL" toml:"level"`
}

type DBConf struct {
	Host     string `mapstructure:"DBHOST" toml:"host"`
	Port     int    `mapstructure:"DBPORT" toml:"port"`
	DBName   string `mapstructure:"DBNAME" toml:"port"`
	User     string `mapstructure:"DBUSER" toml:"user"`
	Password string `mapstructure:"DBPASSWORD" toml:"password"`
}

func (db *DBConf) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/", db.User, db.Password, db.Host, db.Port)
}

type APIConf struct {
	Host string `mapstructure:"APIHOST" toml:"host"`
	Port int    `mapstructure:"APIPORT" toml:"port"`
}

func (cfg *Config) bindEnv() {
	viper.AutomaticEnv()
	viper.Unmarshal(&cfg.Logger)
	viper.Unmarshal(&cfg.DB)
	viper.Unmarshal(&cfg.API)
}

func (cfg *Config) bindFile(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			viper.Unmarshal(cfg)
		}
	}
}

func New(cfgFile string) Config {
	cfg := Config{
		Logger: LoggerConf{Level: "debug"},
		DB:     DBConf{},
		API:    APIConf{Port: 8083},
	}

	cfg.bindEnv()
	cfg.bindFile(cfgFile)

	return cfg
}
