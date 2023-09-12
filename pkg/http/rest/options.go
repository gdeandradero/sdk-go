package rest

import (
	"net/http"
	"time"
)

type options struct {
	timeout       time.Duration
	customHeaders http.Header
}

type Option interface {
	apply(*options)
}

type timeoutOption time.Duration

func (t timeoutOption) apply(opts *options) {
	opts.timeout = time.Duration(t)
}

func WithTimeout(t time.Duration) Option {
	return timeoutOption(t)
}

type customHeadersOption http.Header

func (c customHeadersOption) apply(opts *options) {
	opts.customHeaders = http.Header(c)
}

func WithCustomHeaders(h http.Header) Option {
	return customHeadersOption(h)
}
