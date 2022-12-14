package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"wckd1/tg-youtube-podcasts-bot/internal/bot"
	db "wckd1/tg-youtube-podcasts-bot/internal/db"
	"wckd1/tg-youtube-podcasts-bot/internal/feed"
	"wckd1/tg-youtube-podcasts-bot/internal/file_manager"
	"wckd1/tg-youtube-podcasts-bot/internal/handlers"
	"wckd1/tg-youtube-podcasts-bot/internal/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	bolt "go.etcd.io/bbolt"
)

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal("[ERROR] cannot load config:", err)
	}

	dbConn, err := bolt.Open("storage/yt_podcasts.db", 0666, nil)
	if err != nil {
		log.Fatal("[ERROR] cannot connect to database:", err)
	}
	dbStore := db.NewStore(dbConn)

	ctx, cancel := context.WithCancel(context.Background())
	tgAPI, err := tgbotapi.NewBotAPI(config.Telegram.BotAPIToken)
	if err != nil {
		log.Fatal("[ERROR] cannot init bot api:", err)
	}

	tgAPI.Debug = config.Telegram.DebugMode

	// Grasefull quit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		dbConn.Close()
		cancel()
	}()

	// File manager
	fileManager := file_manager.FileManager{
		Downloader: &file_manager.YTDLPLoader{},
		Uploader: &file_manager.TelegramUploader{
			BotAPI: tgAPI,
			ChatID: config.Telegram.ChatID,
		},
	}

	// Feed service
	feedSrv := feed.FeedService{
		Context:     ctx,
		Limit:       config.Feed.Limit,
		Store:       dbStore,
		FileManager: fileManager,
	}

	// Config available commands
	commands := bot.Commands{
		bot.Add{FeedService: feedSrv},
		bot.Remove{FeedService: feedSrv},
	}

	// Telegram listener for handle commands
	tgListener := handlers.Telegram{
		BotAPI:   tgAPI,
		Commands: commands,
		ChatID:   config.Telegram.ChatID,
	}

	// Server for API
	server := handlers.Server{
		FeedService: feedSrv,
	}

	// Timer handler for handle updates
	updater := handlers.Updater{
		Delay:       config.Feed.UpdateInterval,
		FeedService: feedSrv,
	}

	// Start handlers
	go updater.Start(ctx)
	go server.Run(ctx, config.Server.RssKey, config.Server.Port)

	if err = tgListener.Start(ctx); err != nil {
		log.Printf("[INFO] telegram listener stopped, %v", err)
	}
}
