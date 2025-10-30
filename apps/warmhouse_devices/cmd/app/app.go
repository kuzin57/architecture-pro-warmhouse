package app

import (
	"context"
	"log"
	"os"

	gatesapi "github.com/warmhouse/warmhouse_devices/internal/clients/gates_api"
	temperatureapi "github.com/warmhouse/warmhouse_devices/internal/clients/temperature_api"
	"github.com/warmhouse/warmhouse_devices/internal/config"
	"github.com/warmhouse/warmhouse_devices/internal/cron"
	"github.com/warmhouse/warmhouse_devices/internal/generated/async"
	"github.com/warmhouse/warmhouse_devices/internal/repositories"
	devicesrepo "github.com/warmhouse/warmhouse_devices/internal/repositories/devices"
	telemetryrepo "github.com/warmhouse/warmhouse_devices/internal/repositories/telemetry"
	"github.com/warmhouse/warmhouse_devices/internal/router"
	"github.com/warmhouse/warmhouse_devices/internal/services/common"
	"github.com/warmhouse/warmhouse_devices/internal/services/gates"
	"github.com/warmhouse/warmhouse_devices/internal/services/warming"

	"github.com/warmhouse/libraries/rabbitmq"
	"github.com/warmhouse/libraries/scheduler"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func mustLoadConfig(confPath string) *config.Config {
	data, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	conf := &config.Config{}

	if err := yaml.Unmarshal(data, conf); err != nil {
		panic(err)
	}

	return conf
}

func mustLoadSecrets(secretsPath string) *config.Secrets {
	data, err := os.ReadFile(secretsPath)
	if err != nil {
		log.Fatalf("failed to load secrets: %v", err)
	}

	secrets := &config.Secrets{}

	if err := yaml.Unmarshal(data, secrets); err != nil {
		panic(err)
	}

	return secrets
}

func Create(confPath, secretsPath string) fx.Option {
	var (
		conf    = mustLoadConfig(confPath)
		secrets = mustLoadSecrets(secretsPath)
	)

	return fx.Options(
		fx.Supply(conf, secrets),
		fx.Supply(&secrets.RabbitMQ),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.RecoverFromPanics(),
		fx.Provide(
			zap.NewProduction,
			repositories.NewPgDriver,
			fx.Annotate(devicesrepo.NewRepository, fx.As(new(warming.DevicesRepository))),
			fx.Annotate(devicesrepo.NewRepository, fx.As(new(gates.DevicesRepository))),
			fx.Annotate(devicesrepo.NewRepository, fx.As(new(cron.DevicesRepository))),
			fx.Annotate(telemetryrepo.NewRepository, fx.As(new(cron.TelemetryRepository))),
			telemetryrepo.NewRepository,
			common.NewStarter,
			devicesrepo.NewRepository,
			temperatureapi.NewClient,
			gatesapi.NewClient,
			fx.Annotate(gates.NewService, fx.As(new(router.GatesService))),
			fx.Annotate(warming.NewService, fx.As(new(router.WarmingService))),
			fx.Annotate(gatesapi.NewClient, fx.As(new(gates.GatesClient))),
			fx.Annotate(gatesapi.NewClient, fx.As(new(cron.GatesClient))),
			fx.Annotate(temperatureapi.NewClient, fx.As(new(warming.TemperatureClient))),
			fx.Annotate(temperatureapi.NewClient, fx.As(new(cron.TemperatureClient))),
			fx.Annotate(rabbitmq.NewRabbitMQBrokerController, fx.As(new(extensions.BrokerController))),
			scheduler.NewScheduler,
			async.NewAppController,
			router.NewRouter,
			NewApp,
		),
		fx.Invoke(func(app *App) {}),
		fx.Invoke(func(scheduler *scheduler.Scheduler) {}),
		fx.Invoke(func(starter *common.Starter) {}),
	)
}

type App struct {
	router *router.Router
	broker *async.AppController
	ctx    context.Context
	cancel context.CancelFunc
}

func NewApp(lc fx.Lifecycle, router *router.Router, broker *async.AppController) *App {
	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		router: router,
		broker: broker,
		ctx:    ctx,
		cancel: cancel,
	}

	lc.Append(fx.Hook{
		OnStart: app.Run,
		OnStop:  app.Stop,
	})

	return app
}

func (a *App) Run(ctx context.Context) error {
	if err := a.broker.SubscribeToAllChannels(a.ctx, a.router); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.cancel()
	a.broker.UnsubscribeFromAllChannels(ctx)

	return nil
}
