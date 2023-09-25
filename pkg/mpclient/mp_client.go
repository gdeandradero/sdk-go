package mpclient

import (
	"io"
	"net/http"

	"github.com/gdeandradero/sdk-go/pkg/mpclient/rest"
)

type MercadoPago interface {
	SendRest(reqConfig RequestConfig, opts ...rest.Option) ([]byte, error)
}

type client struct {
	rc rest.Client
}

type RequestConfig struct {
	Method string
	URL    string

	Body io.Reader
}

func New() MercadoPago {
	return &client{rc: rest.Instance()}
}

func (c *client) SendRest(reqConfig RequestConfig, opts ...rest.Option) ([]byte, error) {
	req, err := http.NewRequest(reqConfig.Method, reqConfig.URL, reqConfig.Body)
	if err != nil {
		return nil, &rest.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error creating request: " + err.Error(),
		}
	}

	res, err := c.rc.Send(req, opts...)
	if err != nil {
		return nil, &rest.ErrorResponse{
			StatusCode: res.StatusCode,
			Message:    "error sending request: " + err.Error(),
		}
	}
	defer res.Body.Close()

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &rest.ErrorResponse{
			StatusCode: res.StatusCode,
			Message:    "error reading response body: " + err.Error(),
			Headers:    res.Header,
		}
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, &rest.ErrorResponse{
			StatusCode: res.StatusCode,
			Message:    string(response),
			Headers:    res.Header,
		}
	}

	return response, nil
}
