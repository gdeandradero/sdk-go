package payment

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gdeandradero/sdk-go/pkg/mpclient"
	"github.com/gdeandradero/sdk-go/pkg/mpclient/rest"
)

const (
	postURL   = "https://api.mercadopago.com/v1/payments"
	searchURL = "https://api.mercadopago.com/v1/payments/search"
	getURL    = "https://api.mercadopago.com/v1/payments/{id}"
	putURL    = "https://api.mercadopago.com/v1/payments/{id}"
)

// Client contains the methods to interact with the Payments API.
type Client interface {
	/*
		Create creates a new payment.
		It is a post request to the endpoint: https://api.mercadopago.com/v1/payments
		Reference: https://www.mercadopago.com.br/developers/pt/reference/payments/_payments/post/
	*/
	Create(dto Request, opts ...rest.Option) (*Response, error)

	/*
		Search searches for payments.
		It is a get request to the endpoint: https://api.mercadopago.com/v1/payments/search
		Reference: https://www.mercadopago.com.br/developers/pt/reference/payments/_payments_search/get/
	*/
	Search(f Filters, opts ...rest.Option) (*SearchResponse, error)

	/*
		Get gets a payment by its ID.
		It is a get request to the endpoint: https://api.mercadopago.com/v1/payments/{id}
		Reference: https://www.mercadopago.com.br/developers/pt/reference/payments/_payments_id/get/
	*/
	Get(id int64, opts ...rest.Option) (*Response, error)

	/*
		Cancel cancels a payment by its ID.
		It is a put request to the endpoint: https://api.mercadopago.com/v1/payments/{id}
	*/
	Cancel(id int64, opts ...rest.Option) (*Response, error)

	/*
		Capture captures a payment by its ID.
		It is a put request to the endpoint: https://api.mercadopago.com/v1/payments/{id}
	*/
	Capture(id int64, opts ...rest.Option) (*Response, error)

	/*
		CaptureAmount captures amount of a payment by its ID.
		It is a put request to the endpoint: https://api.mercadopago.com/v1/payments/{id}
	*/
	CaptureAmount(id int64, amount float64, opts ...rest.Option) (*Response, error)
}

// client is the implementation of Client.
type client struct {
	mpc mpclient.MercadoPago
}

// NewClient returns a new Payments API Client.
func NewClient() Client {
	return &client{mpc: mpclient.New()}
}

func (c *client) Create(dto Request, opts ...rest.Option) (*Response, error) {
	body, err := json.Marshal(&dto)
	if err != nil {
		return nil, &rest.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error marshaling request body: " + err.Error(),
		}
	}

	reqConfig := mpclient.RequestConfig{
		Method: http.MethodPost,
		URL:    postURL,
		Body:   strings.NewReader(string(body)),
	}

	res, err := c.mpc.SendRest(reqConfig, opts...)
	if err != nil {
		return nil, err
	}

	formatted := &Response{}
	if err := json.Unmarshal(res, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}

func (c *client) Search(f Filters, opts ...rest.Option) (*SearchResponse, error) {
	params := url.Values{}
	params.Add("sort", f.Sort)
	params.Add("criteria", f.Criteria)
	params.Add("external_reference", f.ExternalReference)
	params.Add("range", f.Range)
	params.Add("begin_date", f.BeginDate)
	params.Add("end_date", f.EndDate)

	reqConfig := mpclient.RequestConfig{
		Method: http.MethodGet,
		URL:    searchURL + "?" + params.Encode(),
		Body:   nil,
	}

	res, err := c.mpc.SendRest(reqConfig, opts...)
	if err != nil {
		return nil, err
	}

	var formatted *SearchResponse
	if err := json.Unmarshal(res, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}

func (c *client) Get(id int64, opts ...rest.Option) (*Response, error) {
	conv := strconv.Itoa(int(id))

	reqConfig := mpclient.RequestConfig{
		Method: http.MethodGet,
		URL:    strings.Replace(getURL, "{id}", conv, 1),
		Body:   nil,
	}

	res, err := c.mpc.SendRest(reqConfig, opts...)
	if err != nil {
		return nil, err
	}

	formatted := &Response{}
	if err := json.Unmarshal(res, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}

func (c *client) Cancel(id int64, opts ...rest.Option) (*Response, error) {
	dto := &CancelRequest{Status: "cancelled"}
	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	conv := strconv.Itoa(int(id))
	reqConfig := mpclient.RequestConfig{
		Method: http.MethodPut,
		URL:    strings.Replace(putURL, "{id}", conv, 1),
		Body:   strings.NewReader(string(body)),
	}

	res, err := c.mpc.SendRest(reqConfig, opts...)
	if err != nil {
		return nil, err
	}

	formatted := &Response{}
	if err := json.Unmarshal(res, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}

func (c *client) Capture(id int64, opts ...rest.Option) (*Response, error) {
	dto := &CaptureRequest{Capture: true}
	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	conv := strconv.Itoa(int(id))
	reqConfig := mpclient.RequestConfig{
		Method: http.MethodPut,
		URL:    strings.Replace(putURL, "{id}", conv, 1),
		Body:   strings.NewReader(string(body)),
	}

	res, err := c.mpc.SendRest(reqConfig, opts...)
	if err != nil {
		return nil, err
	}

	formatted := &Response{}
	if err := json.Unmarshal(res, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}

func (c *client) CaptureAmount(id int64, amount float64, opts ...rest.Option) (*Response, error) {
	dto := &CaptureRequest{TransactionAmount: amount, Capture: true}
	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	conv := strconv.Itoa(int(id))
	reqConfig := mpclient.RequestConfig{
		Method: http.MethodPut,
		URL:    strings.Replace(putURL, "{id}", conv, 1),
		Body:   strings.NewReader(string(body)),
	}

	res, err := c.mpc.SendRest(reqConfig, opts...)
	if err != nil {
		return nil, err
	}

	formatted := &Response{}
	if err := json.Unmarshal(res, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}
