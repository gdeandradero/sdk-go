package payment

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gdeandradero/sdk-go/pkg/http/rest"
)

const (
	postURL   = "https://api.mercadopago.com/v1/payments"
	searchURL = "https://api.mercadopago.com/v1/payments/search"
	getURL    = "https://api.mercadopago.com/v1/payments/{id}"
)

// Client contains the methods to interact with the Payments API.
type Client interface {
	/*
		Create creates a new payment.
		It is a post request to the endpoint: https://api.mercadopago.com/v1/payments
		Reference: https://www.mercadopago.com.br/developers/pt/reference/payments/_payments/post/
	*/
	Create(dto Request) (*Response, error)

	/*
		Search searches for payments.
		It is a get request to the endpoint: https://api.mercadopago.com/v1/payments/search
		Reference: https://www.mercadopago.com.br/developers/pt/reference/payments/_payments_search/get/
	*/
	Search(f Filters) (*SearchResponse, error)

	/*
		Get gets a payment by its ID.
		It is a get request to the endpoint: https://api.mercadopago.com/v1/payments/{id}
		Reference: https://www.mercadopago.com.br/developers/pt/reference/payments/_payments_id/get/
	*/
	Get(id int64) (*Response, error)
}

// client is the implementation of Client.
type client struct {
	hc rest.Client
}

// NewClient returns a new Payments API Client.
func NewClient() Client {
	return &client{hc: rest.Instance()}
}

func (c *client) Create(dto Request) (*Response, error) {
	body, err := json.Marshal(&dto)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(string(body))
	req, err := http.NewRequest("POST", postURL, reader)
	if err != nil {
		return nil, err
	}

	res, err := c.hc.Send(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	read, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, &rest.ErrorResponse{
			StatusCode: res.StatusCode,
			Message:    string(read),
			Headers:    res.Header,
		}
	}

	formatted := &Response{}
	if err := json.Unmarshal(read, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}

func (c *client) Search(f Filters) (*SearchResponse, error) {
	params := url.Values{}
	params.Add("sort", f.Sort)
	params.Add("criteria", f.Criteria)
	params.Add("external_reference", f.ExternalReference)
	params.Add("range", f.Range)
	params.Add("begin_date", f.BeginDate)
	params.Add("end_date", f.EndDate)

	req, err := http.NewRequest("GET", searchURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.hc.Send(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	read, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, &rest.ErrorResponse{
			StatusCode: res.StatusCode,
			Message:    string(read),
			Headers:    res.Header,
		}
	}

	var formatted *SearchResponse
	if err := json.Unmarshal(read, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}

func (c *client) Get(id int64) (*Response, error) {
	conv := strconv.Itoa(int(id))
	url := strings.Replace(getURL, "{id}", conv, 1)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.hc.Send(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	read, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, &rest.ErrorResponse{
			StatusCode: res.StatusCode,
			Message:    string(read),
			Headers:    res.Header,
		}
	}

	formatted := &Response{}
	if err := json.Unmarshal(read, &formatted); err != nil {
		return nil, err
	}

	return formatted, nil
}
