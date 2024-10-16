package main

import (
	"fmt"
	"github.com/spf13/viper"
)

// Load config file using viper
func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}
}
