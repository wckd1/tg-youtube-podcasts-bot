package app

import (
	"context"
	"log"
	"net/http"
	"wckd1/tg-youtube-podcasts-bot/configs"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/httpapi"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram/command"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/updater"
)

type App struct {
	ctx    context.Context
	config configs.Config

	serviceProvider  *serviceProvider
	telegramListener *telegram.TelegramListener
	httpServer       *httpapi.HTTPServer
	updater          updater.Updater
}

func NewApp(ctx context.Context, config configs.Config) (*App, error) {
	a := &App{
		ctx:    ctx,
		config: config,
	}

	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initTelegramListener,
		a.initHTTPServer,
		a.initUpdater,
	}
	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (a App) Run() {
	go func() {
		a.telegramListener.Start(a.ctx)
	}()

	go a.updater.Start(a.ctx)

	go func() {
		if err := a.httpServer.Start(a.ctx); err != nil && err != context.Canceled && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] can't start api server: %+v", err)
		}
	}()
}

func (a App) Shutdown(ctx context.Context) error {
	a.telegramListener.Shutdown()
	a.updater.Shutdown()
	return a.httpServer.Shutdown(ctx)
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.ctx, a.config)
	return nil
}

func (a *App) initTelegramListener(_ context.Context) error {
	listener, err := telegram.NewTelegramListener(
		a.config.Telegram.BotAPIToken,
		a.config.Telegram.DebugMode,
	)
	if err != nil {
		return err
	}

	listener.RegisterCommands(
		command.NewRegisterCommand(a.serviceProvider.RegisterUsecase()),
		command.NewAddCommand(a.serviceProvider.AddUsecase()),
		command.NewPlaylistCommand(*a.serviceProvider.PlaylistUsecase()),
	)

	a.telegramListener = listener
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	a.httpServer = httpapi.NewServer(a.config.Server.Port)

	a.httpServer.RegisterRSSHandler(a.serviceProvider.RSSUseCase())

	return nil
}

func (a *App) initUpdater(_ context.Context) error {
	a.updater = updater.NewUpdater(
		a.serviceProvider.UpdateUsecase(),
		a.serviceProvider.ContentManager(),
		a.config.Feed.UpdateInterval,
	)
	return nil
}
