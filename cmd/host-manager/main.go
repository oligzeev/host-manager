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

func main() {
	ctx, done := context.WithCancel(context.Background())
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return util.StartSignalReceiver(groupCtx, done)
	})

	cfg := initConfig("config/host-manager.yaml", "app_")

	// TODO add linter
	// TODO add metrics
	// TODO add profiler
	// TODO add tests
	// TODO add benchmarks
	// TODO read config prefix from env
	// TODO add openshift's routes watcher
	initLogger(cfg.Logging)

	mappingService := initMappingService(cfg.Mapping)
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
	if err := group.Wait(); err != nil && err != context.Canceled {
		log.Error(err)
	}
}

func initMappingService(cfg domain.MappingConfig) domain.MappingService {
	envMappingService := service.NewEnvMappingService(cfg)
	openshiftMappingService, err := service.NewOpenshiftMappingService(cfg)
	if err != nil {
		log.Fatal(err)
	}
	mappingService := service.NewAggMappingService(envMappingService, openshiftMappingService)
	return mappingService
}

func initConfig(yamlFileName, envPrefix string) *domain.ApplicationConfig {
	cfg, err := util.ReadConfig(yamlFileName, envPrefix)
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
