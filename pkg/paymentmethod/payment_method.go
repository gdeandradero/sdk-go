package paymentmethod

import (
	"encoding/json"
	"net/http"

	"github.com/gdeandradero/sdk-go/pkg/mpclient"
	"github.com/gdeandradero/sdk-go/pkg/mpclient/rest"
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
	mpc mpclient.MercadoPago
}

// NewClient returns a new Payment Methods API Client.
func NewClient() Client {
	return &client{mpc: mpclient.New()}
}

func (c *client) List(opts ...rest.Option) ([]Response, error) {
	reqConfig := mpclient.RequestConfig{
		Method: http.MethodGet,
		URL:    url,
		Body:   nil,
	}

	res, err := c.mpc.SendRest(reqConfig, opts...)
	if err != nil {
		return nil, err
	}

	var formatted []Response
	if err := json.Unmarshal(res, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}
