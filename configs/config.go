package configs

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

var (
	ErrParseConfig = errors.New("can't parse config")
	ErrDecode      = errors.New("unable to decode into struct")
)

type Config struct {
	Feed struct {
		UpdateInterval time.Duration `mapstructure:"update_interval"`
	}
	Server struct {
		Port string `mapstructure:"port"`
	}
	Telegram struct {
		BotAPIToken string `mapstructure:"bot_token"`
		DebugMode   bool   `mapstructure:"debug_mode"`
	}
}

func LoadConfig() (Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, errors.Join(ErrParseConfig, err)
	}

	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		return Config{}, errors.Join(ErrDecode, err)
	}

	return c, nil
}
