package main

import "github.com/spf13/viper"


type Config struct {
    dsn string
}

func init() {
    viper.AutomaticEnv()

    viper.SetEnvPrefix("API_PLANT")
    viper.SetDefault("dsn", "postgres://postgres:postgres@postgresql/plants")
}

func NewConfig() *Config {
    return &Config{
        dsn: viper.GetString("dsn"),
    }
}
