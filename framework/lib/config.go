package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func setDefaultConfig() {
	viper.SetDefault("addr", ":50051")
	viper.SetDefault("debug", false)
	viper.SetDefault("kafka-addr", "kafka:9093")
}

func Load() {
	_, file, _, _ := runtime.Caller(2)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Dir(file))

	err := viper.ReadInConfig()

	if err != nil {
		log.Printf("config file read error: %s \n", err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	viper.AutomaticEnv()

	setDefaultConfig()

	viper.WatchConfig()
}

func Get[T any](key string) T {
	r, ok := viper.Get(key).(T)

	if !ok {
		log.Fatalf("key %s has wrong type", key)
	}

	return r
}
