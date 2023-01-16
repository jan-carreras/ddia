package server

import (
	"ddia/src/logger"
)

type options struct {
	logger            logger.Logger
	dbs               []Storage
	host              string
	port              int
	password          string
	configurationFile string
}

// Option defines an interface that all options must match
type Option interface {
	apply(*options)
}

type loggerOption struct {
	log logger.Logger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.log
}

// WithLogger uses log as logger
func WithLogger(log logger.Logger) Option {
	return loggerOption{log: log}
}

type dbs []Storage

func (d dbs) apply(opts *options) {
	opts.dbs = d
}

// WithDBs sets the databases
func WithDBs(db []Storage) Option {
	return dbs(db)

}

type configurationFile string

func (c configurationFile) apply(opts *options) {
	opts.configurationFile = string(c)
}

// WithConfigurationFile defines the path (relative or absolute) where a configuration file for the server is located
func WithConfigurationFile(config string) Option {
	return configurationFile(config)
}

type host string

func (h host) apply(opts *options) {
	opts.host = string(h)
}

// WithHost uses h hostname to start the Redis server
func WithHost(h string) Option {
	return host(h)
}

type port int

func (p port) apply(opts *options) {
	opts.port = int(p)
}

// WithPort starts Redis server on the given port
func WithPort(p int) Option {
	return port(p)
}

// WithRandomPort starts Redis server on a free random port
func WithRandomPort() Option {
	return port(0)
}
