package main

import (
	"fmt"
	"net/http"

	"github.com/gdeandradero/sdk-go/pkg/mp"
	"github.com/gdeandradero/sdk-go/pkg/paymentmethod"
)

type customRetryClient struct{}

func (*customRetryClient) Retry(req *http.Request, httpClient *http.Client, opts ...mp.Option) (*http.Response, error) {
	// some retry implementation
	return nil, nil
}

func main() {
	mp.SetAccessToken("TEST-640110472259637-071923-a761f639c4eb1f0835ff7611f3248628-793910800")

	customRetryClient := &customRetryClient{}
	mp.SetRetryClient(customRetryClient)

	pmc := paymentmethod.NewClient()
	res, err := pmc.List()
	if err != nil {
		panic(err)
	}

	for _, v := range res {
		fmt.Println(v)
	}
}
