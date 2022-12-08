package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Feed struct {
		UpdateInterval time.Duration `mapstructure:"update_interval"`
		Limit          int           `mapstructure:"limit"`
	}
	Server struct {
		RssKey string `mapstructure:"rss_key"`
		Port   int    `mapstructure:"port"`
	}
	Telegram struct {
		BotAPIToken string `mapstructure:"bot_token"`
		ChatID      int64  `mapstructure:"chat_id"`
		DebugMode   bool   `mapstructure:"debug_mode"`
	}
	AWS AWSConfig
}

type AWSConfig struct {
	Key string `mapstructure:"access_key_id"`
	Secret string `mapstructure:"secret_access_key"`
	Region string `mapstructure:"region"`
	Bucket string `mapstructure:"bucket"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("[ERROR] error parse config: %v", err)
		return
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Printf("[ERROR] unable to decode into struct, %v", err)
	}

	return
}
