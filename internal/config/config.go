// Package config provides the configuration for the application. It contains the configuration schema and the default values for the configuration.
package config

import "time"

// Config represents the configuration for the application.
type Config struct {
	ConfigPath  string      `mapstructure:"config_path"`
	LogLevel    string      `mapstructure:"log_level"`
	Stacktrace  bool        `mapstructure:"stacktrace"`
	Placeholder Placeholder `mapstructure:"placeholder"`
	Server      Server      `mapstructure:"server"`
	Database    Database    `mapstructure:"database"`
}

// Placeholder represents the configuration for the Placeholder command.
type Placeholder struct {
	ID int `mapstructure:"id"`
}

// Server holds the configuration for the server.
type Server struct {
	Host    string        `mapstructure:"host"`
	Port    int           `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// Database holds database configuration.
type Database struct {
	Driver            string        `mapstructure:"driver"`
	DatabaseURL       string        `mapstructure:"database_url"`
	MaxConnection     int           `mapstructure:"max_connection"`
	MaxIdleConnection int           `mapstructure:"max_idle_connection"`
	ConnMaxLifetime   time.Duration `mapstructure:"conn_max_lifetime"`
	PingTimeout       time.Duration `mapstructure:"ping_timeout"`
}
