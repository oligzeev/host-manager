package main

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	"github.com/oligzeev/host-manager/internal/logging"
	"github.com/oligzeev/host-manager/internal/rest"
	"github.com/oligzeev/host-manager/internal/service"
	"github.com/oligzeev/host-manager/internal/util"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
	"os"

	_ "github.com/oligzeev/host-manager/api/swagger"
)

const (
	envConfigPath = "ENV_CONFIG_PATH"
	envPrefix     = "ENV_PREFIX"

	defaultConfigPath = "config/host-manager.yaml"
	defaultPrefix     = "app"
)

func main() {
	// Initialize error group & signal receiver
	ctx, done := context.WithCancel(context.Background())
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return util.StartSignalReceiver(groupCtx, done)
	})

	// Initialize configuration
	cfg := initConfig()

	// TODO add golangci-lint
	// TODO add metrics
	// TODO add profiler
	// TODO add tests
	// TODO add benchmarks
	// TODO add opentracing

	// Initialize logging
	initLogger(cfg.Logging)

	// Initialize mapping services
	envMappingService := service.NewEnvMappingService(cfg.Mapping, os.Environ())
	osMappingService, err := service.NewOpenshiftMappingService(cfg.Mapping)
	if err != nil {
		log.Fatal(err)
	}
	group.Go(func() error {
		osMappingService.StartInformer()
		return nil
	})
	group.Go(func() error {
		<-ctx.Done()
		osMappingService.StopInformer()
		return nil
	})
	mappingService := service.NewAggMappingService(envMappingService, osMappingService)

	// Initialize rest server
	restServer := initRestServer(cfg.Rest.Server, []domain.RestHandler{
		rest.NewMappingRestHandler(mappingService),
	})
	group.Go(func() error {
		return restServer.Start(groupCtx)
	})
	group.Go(func() error {
		<-ctx.Done()
		return restServer.Stop(groupCtx)
	})

	// ...
	if err := group.Wait(); err != nil && err != context.Canceled {
		log.Error(err)
	}
}

func initConfig() *domain.ApplicationConfig {
	configPath := util.GetEnv(envConfigPath, defaultConfigPath)
	prefix := util.GetEnv(envPrefix, defaultPrefix)
	cfg, err := util.ReadConfig(configPath, prefix)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func initLogger(cfg domain.LoggingConfig) {
	if cfg.Default {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		log.SetFormatter(&logging.TextFormatter{
			TimestampFormat: cfg.TimestampFormat,
		})
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.Level(cfg.Level))
}

func initRestServer(cfg domain.ServerRestConfig, handlers []domain.RestHandler) *rest.Server {
	restServer := rest.NewServer(cfg, handlers)
	router := restServer.Router()

	// Swagger handler initialization
	router.GET(cfg.SwaggerUrl+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return restServer
}
