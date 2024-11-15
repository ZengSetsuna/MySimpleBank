package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `default:"postgres"`
	DBSource      string `default:"postgresql://shu:shu@localhost:5432/simple_bank?sslmode=disable"`
	ServerAddress string `default:":8080"`
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
	if err != nil {
		return
	}
	return
}
