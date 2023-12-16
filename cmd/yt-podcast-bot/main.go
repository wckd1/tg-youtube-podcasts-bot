package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wckd1/tg-youtube-podcasts-bot/configs"
	"wckd1/tg-youtube-podcasts-bot/internal/app"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("[ERROR] %+v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	a, err := app.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("[ERROR] %+v", err)
	}

	a.Run()

	<-quit
	log.Println("[INFO] shutting app down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownCancel()

	// Do finishing stuff
	if err := a.Shutdown(shutdownCtx); err != nil {
		log.Printf("[INFO] failed to shut down app, %+v", err)
	}

	<-shutdownCtx.Done()

	log.Println("[INFO] app shutted down")
}
