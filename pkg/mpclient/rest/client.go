package rest

import (
	"context"
	"net/http"

	"github.com/gdeandradero/sdk-go/pkg/config"
	"github.com/google/uuid"
)

var (
	instance Client

	authorizationHeader = http.CanonicalHeaderKey("authorization")
	productIDHeader     = http.CanonicalHeaderKey("x-product-id")
	idempotencyHeader   = http.CanonicalHeaderKey("x-idempotency-key")
)

// Client is the interface that wraps the basic Send method.
type Client interface {
	/*
		Send sends a request to the API.
		opts are optional parameters to be used in the request, if you do not need, ignore it.
	*/
	Send(req *http.Request, opts ...Option) (*http.Response, error)
}

// client is the implementation of Client.
type client struct {
	hc *http.Client
}

// Instance returns a current Client instance or create a new one.
func Instance() Client {
	if instance == nil {
		instance = &client{hc: &http.Client{}}
	}
	return instance
}

// SetCustomHTTPClient sets a custom http.Client to be used by the Client.
func SetCustomHTTPClient(chc *http.Client) {
	instance = &client{hc: chc}
}

func (c *client) Send(req *http.Request, opts ...Option) (*http.Response, error) {
	options := &options{}
	for _, opt := range opts {
		opt.apply(options)
	}
	if options.timeout > 0 {
		ctx, cancel := context.WithTimeout(req.Context(), options.timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}
	if options.customHeaders != nil {
		for k, v := range options.customHeaders {
			canonicalKey := http.CanonicalHeaderKey(k)
			req.Header[canonicalKey] = v
		}
	}
	setDefaultHeaders(req)

	return c.hc.Do(req)
}

func setDefaultHeaders(req *http.Request) {
	req.Header.Add(authorizationHeader, "Bearer "+config.AccessToken())
	req.Header.Add(productIDHeader, config.ProductID())

	if _, ok := req.Header[idempotencyHeader]; !ok {
		req.Header.Add(idempotencyHeader, uuid.New().String())
	}
}
