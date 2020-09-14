package domain

import "time"

type ServerRestConfig struct {
	Host               string        `yaml:"host"`
	Port               int           `yaml:"port"`
	SwaggerUrl         string        `yaml:"swaggerUrl"`
	MetricsUrl         string        `yaml:"metricsUrl"`
	ReadTimeoutSec     time.Duration `yaml:"readTimeoutSec"`
	WriteTimeoutSec    time.Duration `yaml:"writeTimeoutSec"`
	ShutdownTimeoutSec time.Duration `yaml:"shutdownTimeoutSec"`
	Release            bool          `yaml:"release"`
}

type RestConfig struct {
	Server ServerRestConfig `yaml:"server"`
}

type LoggingConfig struct {
	Level           int    `yaml:"level"`
	TimestampFormat string `yaml:"timestampFormat"`
	Default         bool   `yaml:"default"`
}

type MappingConfig struct {
	Prefix     string `yaml:"prefix"`
	Namespace  string `yaml:"namespace"`
	ConfigPath string `yaml:"configPath"`
}

type TracingConfig struct {
	ServiceName       string `yaml:"serviceName"`
	CollectorEndpoint string `yaml:"collectorEndpoint"`
}

type ApplicationConfig struct {
	Rest    RestConfig    `yaml:"rest"`
	Logging LoggingConfig `yaml:"logging"`
	Mapping MappingConfig `yaml:"mapping"`
	Tracing TracingConfig `yaml:"tracing"`
}
