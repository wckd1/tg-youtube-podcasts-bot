package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"wckd1/tg-youtube-podcasts-bot/bot"
	"wckd1/tg-youtube-podcasts-bot/handlers"
	"wckd1/tg-youtube-podcasts-bot/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	tgAPI, err := tgbotapi.NewBotAPI(config.BotAPIToken)
	if err != nil {
		log.Fatal("Cannot init bot api:", err)
	}

	tgAPI.Debug = config.DebugMode

	// Grasefull quit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		cancel()
	}()

	// Config available commands
	commands := bot.Commands{
		bot.Add{},
	}

	// Telegram listener to handle commands
	tgListener := handlers.TelegramListener{
		BotAPI:   tgAPI,
		Commands: commands,
		ChatID:   -826459712,
	}

	// Timer to handle check for updates
	updateChecker := handlers.UpdateChecker{
		Delay:     config.UpdateInterval,
		Submitter: &tgListener,
	}

	// Start handlers
	go func() {
		if err := updateChecker.Start(ctx); err != nil {
			log.Printf("[INFO] update checker stopped, %v", err)
		}
	}()
	if err := tgListener.Start(ctx); err != nil {
		log.Printf("[INFO] telegram listener stopped, %v", err)
	}
}
