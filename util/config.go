package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER" default:"postgres"`
	DBSource            string        `mapstructure:"DB_SOURCE" default:"postgresql://shu:shu@localhost:5432/simple_bank?sslmode=disable"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS" default:":8080"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION" default:"15m"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	fmt.Println("config: ", config)
	if err != nil {
		return
	}
	return
}
