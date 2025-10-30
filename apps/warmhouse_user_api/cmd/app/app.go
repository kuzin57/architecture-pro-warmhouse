package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/warmhouse/warmhouse_user_api/internal/config"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/async"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
	"github.com/warmhouse/warmhouse_user_api/internal/handlers"
	"github.com/warmhouse/warmhouse_user_api/internal/repositories"
	devicesrepo "github.com/warmhouse/warmhouse_user_api/internal/repositories/devices"
	housesrepo "github.com/warmhouse/warmhouse_user_api/internal/repositories/houses"
	telemetryrepo "github.com/warmhouse/warmhouse_user_api/internal/repositories/telemetry"
	usersrepo "github.com/warmhouse/warmhouse_user_api/internal/repositories/users"
	"github.com/warmhouse/warmhouse_user_api/internal/services/devices"
	"github.com/warmhouse/warmhouse_user_api/internal/services/houses"
	"github.com/warmhouse/warmhouse_user_api/internal/services/users"

	"github.com/warmhouse/libraries/rabbitmq"

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

	var (
		handlerOptions = server.ChiServerOptions{
			BaseURL: "/api/v1",
		}
		strictMiddlewares = []server.StrictMiddlewareFunc{}
	)

	return fx.Options(
		fx.Supply(conf, secrets, handlerOptions, strictMiddlewares, &secrets.RabbitMQ),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			zap.NewProduction,
			handlers.NewRegisterUserHandler,
			handlers.NewLoginUserHandler,
			handlers.NewGetDefaultUserHandler,
			handlers.NewGetUserInfoHandler,
			handlers.NewUpdateUserHandler,
			handlers.NewGetUserHousesHandler,
			handlers.NewCreateHouseHandler,
			handlers.NewDeleteHouseHandler,
			handlers.NewGetHouseInfoHandler,
			handlers.NewUpdateHouseHandler,
			handlers.NewGetUserDevicesHandler,
			handlers.NewCreateDeviceHandler,
			handlers.NewDeleteDeviceHandler,
			handlers.NewGetDeviceInfoHandler,
			handlers.NewUpdateDeviceHandler,
			handlers.NewGetDeviceTelemetryHandler,
			fx.Annotate(handlers.NewServer, fx.As(new(server.StrictServerInterface))),
			fx.Annotate(server.NewStrictHandler, fx.As(new(server.ServerInterface))),
			fx.Annotate(server.HandlerWithOptions, fx.As(new(http.Handler))),
			repositories.NewPgDriver,
			fx.Annotate(usersrepo.NewRepository, fx.As(new(users.UsersRepository))),
			fx.Annotate(housesrepo.NewRepository, fx.As(new(users.HousesRepository))),
			fx.Annotate(users.NewService, fx.As(new(handlers.UsersService))),
			fx.Annotate(housesrepo.NewRepository, fx.As(new(houses.HousesRepository))),
			fx.Annotate(housesrepo.NewRepository, fx.As(new(devices.HouseRepository))),
			fx.Annotate(houses.NewService, fx.As(new(handlers.HousesService))),
			fx.Annotate(devicesrepo.NewRepository, fx.As(new(devices.DevicesRepository))),
			fx.Annotate(devicesrepo.NewRepository, fx.As(new(houses.DevicesRepository))),
			fx.Annotate(telemetryrepo.NewRepository, fx.As(new(devices.TelemetryRepository))),
			fx.Annotate(devices.NewService, fx.As(new(handlers.DevicesService))),
			fx.Annotate(rabbitmq.NewRabbitMQBrokerController, fx.As(new(extensions.BrokerController))),
			async.NewUserController,
			NewHTTPServer,
		),
		fx.Invoke(func(server *HTTPServer) {}),
	)
}

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(lc fx.Lifecycle, conf *config.Config, log *zap.Logger, handler http.Handler) *HTTPServer {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", conf.App.Port),
		Handler: handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server", zap.String("addr", server.Addr))

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("Failed to start HTTP server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down HTTP server")
			return server.Shutdown(ctx)
		},
	})

	return &HTTPServer{server: server}
}
