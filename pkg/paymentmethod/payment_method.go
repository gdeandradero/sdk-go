package paymentmethod

import (
	"encoding/json"
	"net/http"

	"github.com/gdeandradero/sdk-go/pkg/mp"
)

const url = "https://api.mercadopago.com/v1/payment_methods"

// Client contains the methods to interact with the Payment Methods API.
type Client interface {
	/*
		List lists all payment methods.
		It is a get request to the endpoint: https://api.mercadopago.com/v1/payment_methods
		Reference: https://www.mercadopago.com.br/developers/pt/reference/payment_methods/_payment_methods/get/
	*/
	List(opts ...mp.Option) ([]Response, error)
}

// client is the implementation of Client.
type client struct{}

// NewClient returns a new Payment Methods API Client.
func NewClient() Client {
	return &client{}
}

func (c *client) List(opts ...mp.Option) ([]Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, &mp.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error creating request: " + err.Error(),
		}
	}

	res, err := mp.GetRestClient().Send(req, opts...)
	if err != nil {
		return nil, err
	}

	var formatted []Response
	if err := json.Unmarshal(res, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}
