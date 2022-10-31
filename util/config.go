package util

import "github.com/spf13/viper"

type Config struct {
	BotAPIToken string `mapstructure:"BOT_API_TOKEN"`
	DebugMode   bool   `mapstructure:"DEBUG_MODE"`
}

func LoadConfig(path string) (c Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}
