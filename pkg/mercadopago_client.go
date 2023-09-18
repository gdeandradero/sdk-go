package pkg

import (
	"io"
	"net/http"

	"github.com/gdeandradero/sdk-go/pkg/http/rest"
)

type MercadoPagoClient interface {
	Send(req *http.Request, opts ...rest.Option) ([]byte, error)
}

type client struct {
	hc rest.Client
}

func NewMercadoPagoClient() MercadoPagoClient {
	return &client{hc: rest.Instance()}
}

func (c *client) Send(req *http.Request, opts ...rest.Option) ([]byte, error) {
	res, err := c.hc.Send(req, opts...)
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
