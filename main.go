package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"wckd1/tg-youtube-podcasts-bot/bot"
	db "wckd1/tg-youtube-podcasts-bot/db/store"
	"wckd1/tg-youtube-podcasts-bot/handlers"
	"wckd1/tg-youtube-podcasts-bot/util"

	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed db/migration
var migrations embed.FS

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal("[ERROR] cannot load config:", err)
	}

	dbConn, err := sql.Open("sqlite3", "./storage/yt_podcasts.db")
	if err != nil {
		log.Fatal("[ERROR] cannot connect to database:", err)
	}
	if err := db.Migrate(dbConn, migrations); err != nil {
		log.Printf("[INFO] Database migration failed: %v", err)
	}
	dbStore := db.NewStore(dbConn)

	ctx, cancel := context.WithCancel(context.Background())
	tgAPI, err := tgbotapi.NewBotAPI(config.BotAPIToken)
	if err != nil {
		log.Fatal("[ERROR] cannot init bot api:", err)
	}

	tgAPI.Debug = config.DebugMode

	// Grasefull quit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		dbConn.Close()
		cancel()
	}()

	// Config available commands
	commands := bot.Commands{
		bot.Add{
			Context: ctx,
			Store:   dbStore,
		},
	}

	// Telegram listener for handle commands
	tgListener := handlers.TelegramListener{
		BotAPI:   tgAPI,
		Commands: commands,
		ChatID:   config.ChatID,
	}

	// Timer handler for handle updates
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
