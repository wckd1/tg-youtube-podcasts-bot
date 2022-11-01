package main

import (
	"context"
	"log"
	"wckd1/tg-youtube-podcasts-bot/handlers"
	"wckd1/tg-youtube-podcasts-bot/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	ctx := context.TODO()
	tgAPI, err := tgbotapi.NewBotAPI(config.BotAPIToken)
	if err != nil {
		log.Fatal("Cannot init bot api:", err)
	}

	tgAPI.Debug = config.DebugMode

	// Telegram listener to handle commands
	tgListener := handlers.TelegramListener{
		BotAPI: tgAPI,
	}
	tgListener.Start(ctx)
	if err := tgListener.Start(ctx); err != nil {
		log.Fatalf("[ERROR] telegram listener failed, %v", err)
	}
}
