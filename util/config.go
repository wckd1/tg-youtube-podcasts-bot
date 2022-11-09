package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	BotAPIToken    string        `mapstructure:"BOT_API_TOKEN"`
	ChatID         int64         `mapstructure:"CHAT_ID"`
	DebugMode      bool          `mapstructure:"DEBUG_MODE"`
	UpdateInterval time.Duration `mapstructure:"UPDATE_INTERVAL"`
	RssKey         string        `mapstructure:"RSS_KEY"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./")
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
