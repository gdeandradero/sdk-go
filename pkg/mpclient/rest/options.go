package rest

import (
	"net/http"
	"time"
)

type options struct {
	retryQuantity int

	retryWait     time.Duration
	timeout       time.Duration
	customHeaders http.Header
}

type Option interface {
	apply(*options)
}

type retryQuantityOption int

func (rq retryQuantityOption) apply(opts *options) {
	opts.retryQuantity = int(rq)
}

func WithRetryQuantity(q int) Option {
	return retryQuantityOption(q)
}

type retryWaitOption time.Duration

func (rw retryWaitOption) apply(opts *options) {
	opts.retryWait = time.Duration(rw)
}

func WithRetryWait(t time.Duration) Option {
	return retryWaitOption(t)
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
