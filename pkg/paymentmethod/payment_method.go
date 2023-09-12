package paymentmethod

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gdeandradero/sdk-go/pkg/http/rest"
)

const url = "https://api.mercadopago.com/v1/payment_methods"

// Client contains the methods to interact with the Payment Methods API.
type Client interface {
	/*
		List lists all payment methods.
		It is a get request to the endpoint: https://api.mercadopago.com/v1/payment_methods
		Reference: https://www.mercadopago.com.br/developers/pt/reference/payment_methods/_payment_methods/get/
	*/
	List(opts ...rest.Option) ([]Response, error)
}

// client is the implementation of Client.
type client struct {
	hc rest.Client
}

// NewClient returns a new Payment Methods API Client.
func NewClient() Client {
	return &client{hc: rest.Instance()}
}

func (c *client) List(opts ...rest.Option) ([]Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.delegateSend(req, opts...)
	if err != nil {
		return nil, err
	}

	var formatted []Response
	if err := json.Unmarshal(res, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}

func (c *client) delegateSend(req *http.Request, opts ...rest.Option) ([]byte, error) {
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
