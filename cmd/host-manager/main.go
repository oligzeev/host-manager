package main

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	"github.com/oligzeev/host-manager/internal/logging"
	"github.com/oligzeev/host-manager/internal/metric"
	"github.com/oligzeev/host-manager/internal/rest"
	"github.com/oligzeev/host-manager/internal/service/mapping"
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

	"github.com/gin-contrib/pprof"

	_ "github.com/oligzeev/host-manager/api/swagger"
)

const (
	envConfigPath   = "ENV_CONFIG_PATH"
	envConfigPrefix = "ENV_PREFIX"

	defaultConfigPath   = "config/host-manager.yaml"
	defaultConfigPrefix = "app"
)

// Defaults are in config/host-manager.yaml (could be changed via ENV_CONFIG_PATH)
func main() {
	// Initialize error group & signal receiver
	ctx, done := context.WithCancel(context.Background())
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return util.StartSignalReceiver(groupCtx, done)
	})

	// Initialize configuration
	cfg := initConfig()

	// Initialize logging
	initLogger(cfg.Logging)

	// Initialize tracing
	_, closer := initTracing(cfg.Tracing)
	defer closer.Close()

	// Initialize mapping services
	mappingService := initMappingService(cfg.Mapping, func(osMappingService *mapping.OpenshiftMappingService) {
		group.Go(func() error {
			osMappingService.StartInformer()
			return nil
		})
		group.Go(func() error {
			<-ctx.Done()
			osMappingService.StopInformer()
			return nil
		})
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

// Initialize mapping service with tracing and prometheus decorators.
// Values from env-variables and openshift's routes are merged.
func initMappingService(cfg domain.MappingConfig, initInformer func(osMappingService *mapping.OpenshiftMappingService)) domain.MappingService {
	envService := mapping.NewEnvMappingService(cfg, os.Environ())
	if cfg.Namespace != "" {
		osService, err := mapping.NewOpenshiftMappingService(cfg)
		if err != nil {
			log.Fatal(err)
		}
		initInformer(osService)
		return mapping.NewTracingMappingService(
			mapping.NewMetricMappingService(
				mapping.NewAggMappingService(envService, osService)))
	} else {
		return mapping.NewTracingMappingService(
			mapping.NewMetricMappingService(
				mapping.NewAggMappingService(envService, nil)))
	}
}

// Initialize configuration via merging env-variables and config-file.
// ENV_PREFIX - prefix for env-variables
// ENV_CONFIG_PATH - path to config-file
func initConfig() *domain.ApplicationConfig {
	configPath := util.GetEnv(envConfigPath, defaultConfigPath)
	prefix := util.GetEnv(envConfigPrefix, defaultConfigPrefix)
	cfg, err := util.ReadConfig(configPath, prefix)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

// Initialize logrus logger. Default formatter could be changed via configuration
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

// Initialize opentracing and set global tracer
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

// Initialize rest server with application, tracing, swagger, prometheus handlers
func initRestServer(cfg domain.ServerRestConfig, handlers []domain.RestHandler) *rest.Server {
	restServer := rest.NewServer(cfg, handlers)
	router := restServer.Router()

	// Jaeger middleware initialization
	router.Use(tracing.Middleware())

	// Swagger handler initialization
	router.GET(cfg.SwaggerUrl+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prometheus handler initialization
	router.GET(cfg.MetricsUrl, metric.PrometheusHandler())

	// PProf handler initialization
	pprof.Register(router)

	return restServer
}
