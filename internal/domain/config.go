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
	ServiceName string `yaml:"serviceName"`
}

type ApplicationConfig struct {
	Rest    RestConfig    `yaml:"rest"`
	Logging LoggingConfig `yaml:"logging"`
	Mapping MappingConfig `yaml:"mapping"`
	Tracing TracingConfig `yaml:"tracing"`
}
