package main

import (
	"context"
	"log"
	"time"
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

	// Timer to handle check for updates
	updateChecker := handlers.UpdateChecker{
		BotAPI: tgAPI,
	}
	go updateChecker.Start(ctx, time.Second*time.Duration(config.UpdateInterval))

	// Telegram listener to handle commands
	tgListener := handlers.TelegramListener{
		BotAPI: tgAPI,
	}
	tgListener.Start(ctx)
	if err := tgListener.Start(ctx); err != nil {
		log.Fatalf("[ERROR] telegram listener failed, %v", err)
	}
}
