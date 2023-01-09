package server

import (
	"ddia/src/logger"
)

type options struct {
	logger logger.Logger
	host   string
	port   int
	addr   string
}

type Option interface {
	apply(*options)
}

type loggerOption struct {
	log logger.Logger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.log
}

func WithLogger(log logger.Logger) Option {
	return loggerOption{log: log}
}

type host string

func (h host) apply(opts *options) {
	opts.host = string(h)
}

func WithHost(h string) Option {
	return host(h)
}

type port int

func (p port) apply(opts *options) {
	opts.port = int(p)
}

func WithPort(p int) Option {
	return port(p)
}

func WithRandomPort() Option {
	return port(0)
}
