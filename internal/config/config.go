// Package config implements all methods for service configuration.
package config

import (
	"time"
)

type Config struct {
	HTTP         HTTP         `yaml:"http"`
	Postgres     Postgres     `yaml:"postgres"`
	Notification Notification `yaml:"notification"`
}

type HTTP struct {
	Port            string        `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

type Postgres struct {
	URL string `yaml:"url"`
}

type Notification struct {
	BufferSize     int           `yaml:"bufferSize"`
	RecheckTimeout time.Duration `yaml:"recheckTimeout"`
	CloseTimeout   time.Duration `yaml:"closeTimeout"`
}
