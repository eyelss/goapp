package framework

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port int `yaml:"port"`
}

func setDefaultConfig() {
	viper.SetDefault("port", 50051)
	viper.SetDefault("debug", false)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.ReadInConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	viper.AutomaticEnv()

	setDefaultConfig()

	viper.WatchConfig()
}
