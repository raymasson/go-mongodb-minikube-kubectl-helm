package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	DbHost     string
	DbUser     string
	DbPassword string
	DbURI      string
)

//Get : Gets the service configuration
func Get() {
	viper.SetDefault("PORT", "8000")

	if os.Getenv("ENVIRONMENT") == "DEV" {
		_, dirname, _, _ := runtime.Caller(0)
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(filepath.Dir(dirname))
		viper.ReadInConfig()
	} else {
		viper.AutomaticEnv()
	}

	//Assign env variables value to global variables
	DbHost = viper.GetString("DB_HOST")
	DbUser = viper.GetString("DB_USER")
	DbPassword = viper.GetString("DB_PASSWORD")
	DbURI = viper.GetString("DB_URI")
}
