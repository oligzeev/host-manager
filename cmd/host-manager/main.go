package main

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	"github.com/oligzeev/host-manager/internal/logging"
	"github.com/oligzeev/host-manager/internal/metric"
	"github.com/oligzeev/host-manager/internal/rest"
	"github.com/oligzeev/host-manager/internal/service"
	"github.com/oligzeev/host-manager/internal/tracing"
	"github.com/oligzeev/host-manager/internal/util"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uber/jaeger-client-go"
	"golang.org/x/sync/errgroup"
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	jaegerconf "github.com/uber/jaeger-client-go/config"

	_ "github.com/oligzeev/host-manager/api/swagger"
)

const (
	envConfigPath   = "ENV_CONFIG_PATH"
	envConfigPrefix = "ENV_PREFIX"

	defaultConfigPath   = "config/host-manager.yaml"
	defaultConfigPrefix = "app"
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
	// TODO add profiler
	// TODO add tests
	// TODO add benchmarks

	// TODO add metrics
	// TODO add opentracing

	// Initialize logging
	initLogger(cfg.Logging)

	// Initialize tracing
	_, closer := initTracing(cfg.Tracing)
	defer closer.Close()

	// Initialize mapping services
	mappingService, osMappingService := initMappingService(cfg.Mapping)
	group.Go(func() error {
		osMappingService.StartInformer()
		return nil
	})
	group.Go(func() error {
		<-ctx.Done()
		osMappingService.StopInformer()
		return nil
	})

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

func initMappingService(cfg domain.MappingConfig) (domain.MappingService, *service.OpenshiftMappingService) {
	envService := service.NewEnvMappingService(cfg, os.Environ())
	osService, err := service.NewOpenshiftMappingService(cfg)
	if err != nil {
		log.Fatal(err)
	}
	aggService := service.NewAggMappingService(envService, osService)
	metricService := service.NewMetricMappingService(aggService)
	return service.NewTracingMappingService(metricService), osService
}

func initConfig() *domain.ApplicationConfig {
	configPath := util.GetEnv(envConfigPath, defaultConfigPath)
	prefix := util.GetEnv(envConfigPrefix, defaultConfigPrefix)
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

func initTracing(cfg domain.TracingConfig) (opentracing.Tracer, io.Closer) {
	tracingCfg := jaegerconf.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &jaegerconf.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerconf.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := tracingCfg.NewTracer()
	if err != nil {
		log.Fatal(err)
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func initRestServer(cfg domain.ServerRestConfig, handlers []domain.RestHandler) *rest.Server {
	restServer := rest.NewServer(cfg, handlers)
	router := restServer.Router()

	// Jaeger middleware initialization
	router.Use(tracing.Middleware())

	// Swagger handler initialization
	router.GET(cfg.SwaggerUrl+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prometheus handler initialization
	router.GET(cfg.MetricsUrl, metric.PrometheusHandler())

	return restServer
}
