package mp

import "net/http"

// config is an mp instance.
var config = &mp{
	productID: "123",

	restClient:  &restClient{},
	httpClient:  &http.Client{},
	retryClient: &retryClient{},
}

// mp represents the config.
type mp struct {
	accessToken string
	productID   string

	restClient  RestClient
	httpClient  *http.Client
	retryClient RetryClient
}

// SetAccessToken sets the access token.
func SetAccessToken(at string) {
	config.accessToken = at
}

func GetRestClient() RestClient {
	return config.restClient
}

// SetHTTPClient sets a custom http client.
func SetHTTPClient(hc *http.Client) {
	config.httpClient = hc
}

// SetRetryClient sets a custom retry client.
func SetRetryClient(rc RetryClient) {
	config.retryClient = rc
}
