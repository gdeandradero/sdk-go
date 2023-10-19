package mp

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const defaultTimeout = time.Duration(time.Second * 30)

var (
	authorizationHeader = http.CanonicalHeaderKey("authorization")
	productIDHeader     = http.CanonicalHeaderKey("x-product-id")
	idempotencyHeader   = http.CanonicalHeaderKey("x-idempotency-key")
)

// RestClient is the interface that wraps the basic Send method.
type RestClient interface {
	/*
		Send sends a request to the API.
		opts are optional parameters to be used in the request, if you do not need, ignore it.
	*/
	Send(req *http.Request, opts ...Option) ([]byte, error)
}

// client is the implementation of Client.
type restClient struct{}

func (c *restClient) Send(req *http.Request, opts ...Option) ([]byte, error) {
	c.prepareRequest(req, opts...)

	res, err := config.httpClient.Do(req)
	if shouldRetry(res, err) {
		res, err = config.retryClient.Retry(req, config.httpClient, opts...)
	}
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: res.StatusCode,
			Message:    "error sending request: " + err.Error(),
		}
	}

	defer res.Body.Close()

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &ErrorResponse{
			StatusCode: res.StatusCode,
			Message:    "error reading response body: " + err.Error(),
			Headers:    res.Header,
		}
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, &ErrorResponse{
			StatusCode: res.StatusCode,
			Message:    string(response),
			Headers:    res.Header,
		}
	}

	return response, nil
}

func (c *restClient) prepareRequest(req *http.Request, opts ...Option) {
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
	req.Header.Add(authorizationHeader, "Bearer "+config.accessToken)
	req.Header.Add(productIDHeader, config.productID)

	if _, ok := req.Header[idempotencyHeader]; !ok {
		req.Header.Add(idempotencyHeader, uuid.New().String())
	}
}

func shouldRetry(res *http.Response, err error) bool {
	return err != nil || res.StatusCode >= http.StatusInternalServerError
}
