package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gdeandradero/sdk-go/pkg/config"
	"github.com/google/uuid"
)

const defaultTimeout = time.Duration(time.Second * 30)

var (
	currentClient *client

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
	rc RetryClient
}

// Instance returns a current Client instance or create a new one.
func Instance() Client {
	if currentClient == nil {
		currentClient = &client{
			hc: &http.Client{},
			rc: &retry{},
		}
	}
	return currentClient
}

// SetCustomHTTPClient sets a custom http.Client to be used by the Client.
func SetCustomHTTPClient(chc *http.Client) {
	if currentClient == nil {
		_ = Instance()
	}
	currentClient.hc = chc
}

// SetCustomRetryClient sets a custom RetryClient to be used by the Client.
func SetCustomRetryClient(crc RetryClient) {
	if currentClient == nil {
		_ = Instance()
	}
	currentClient.rc = crc
}

func (c *client) Send(req *http.Request, opts ...Option) (*http.Response, error) {
	c.prepareRequest(req, opts...)

	res, err := c.hc.Do(req)
	if shouldRetry(res, err) {
		res, err = c.rc.Retry(req, c.hc, opts...)
	}

	return res, err
}

func (c *client) prepareRequest(req *http.Request, opts ...Option) {
	timeout := defaultTimeout

	options := &options{}
	for _, opt := range opts {
		opt.apply(options)
	}
	if options.timeout > 0 {
		timeout = options.timeout
	}
	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	if options.customHeaders != nil {
		for k, v := range options.customHeaders {
			canonicalKey := http.CanonicalHeaderKey(k)
			req.Header[canonicalKey] = v
		}
	}
	setDefaultHeaders(req)
}

func setDefaultHeaders(req *http.Request) {
	req.Header.Add(authorizationHeader, "Bearer "+config.AccessToken())
	req.Header.Add(productIDHeader, config.ProductID())

	if _, ok := req.Header[idempotencyHeader]; !ok {
		req.Header.Add(idempotencyHeader, uuid.New().String())
	}
}

func shouldRetry(res *http.Response, err error) bool {
	return err != nil || res.StatusCode >= http.StatusInternalServerError
}
