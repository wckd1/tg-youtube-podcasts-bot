package app

import (
	"context"
	"log"
	"wckd1/tg-youtube-podcasts-bot/configs"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/httpapi"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram/command"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/content"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/content/youtube"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/updater"
)

type App struct {
	ctx    context.Context
	config configs.Config

	serviceProvider  *serviceProvider
	telegramListener *telegram.TelegramListener
	httpServer       *httpapi.HTTPServer
	updater          updater.Updater
	contentManager   content.ContentManager
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
		a.initContentManager,
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
		if err := a.telegramListener.Start(a.ctx); err != nil && err != context.Canceled {
			log.Fatalf("[ERROR] can't start telegram listener: %+v", err)
		}
	}()

	go func() {
		if err := a.httpServer.Start(a.ctx); err != nil && err != context.Canceled {
			log.Fatalf("[ERROR] can't start api server: %+v", err)
		}
	}()

	go a.updater.Start(a.ctx)
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
		command.NewAddCommand(a.serviceProvider.SubscriptionUseCase()),
		command.NewRemoveCommand(a.serviceProvider.SubscriptionUseCase()),
	)

	a.telegramListener = listener
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	a.httpServer = httpapi.NewServer(a.config.Server.Port)

	a.httpServer.RegisterRSSHandler(a.serviceProvider.RSSUseCase(), a.config.Server.RssKey)

	return nil
}

func (a *App) initUpdater(_ context.Context) error {
	a.updater = updater.NewUpdater(
		a.serviceProvider.SubscriptionUseCase(),
		a.contentManager,
		a.config.Feed.UpdateInterval,
	)
	return nil
}

func (a *App) initContentManager(_ context.Context) error {
	cm, err := youtube.NewYouTubeContentManager()
	if err != nil {
		return err
	}

	a.contentManager = cm
	return nil
}
