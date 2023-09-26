package rest

import (
	"net/http"
	"time"
)

const (
	defaultRetryQtt = 3
	defaultWait     = time.Second * 5
)

// RetryClient is the interface that defines retry signature.
type RetryClient interface {
	Retry(req *http.Request, httpClient *http.Client, opts ...Option) (*http.Response, error)
}

// retry is the default implementation of RetryClient.
type retry struct{}

func (*retry) Retry(req *http.Request, httpClient *http.Client, opts ...Option) (*http.Response, error) {
	var (
		retryQtt               = defaultRetryQtt
		wait     time.Duration = defaultWait
		res      *http.Response
		err      error
	)

	options := &options{}
	for _, opt := range opts {
		opt.apply(options)
	}
	if options.retryQuantity > 0 {
		retryQtt = options.retryQuantity
	}
	if options.retryWait > 0 {
		wait = options.retryWait
	}

	for i := 0; i < retryQtt; i++ {
		timer := time.NewTimer(wait)
		select {
		case <-req.Context().Done():
			timer.Stop()
			return nil, req.Context().Err()
		case <-timer.C:
		}

		res, err = httpClient.Do(req)
		if shouldStop(res, err) {
			break
		}
	}

	return res, err
}

func shouldStop(res *http.Response, err error) bool {
	return err == nil && res.StatusCode < http.StatusInternalServerError
}
