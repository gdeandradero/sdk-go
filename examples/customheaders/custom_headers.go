package main

import (
	"fmt"
	"net/http"

	"github.com/gdeandradero/sdk-go/pkg/mp"
	"github.com/gdeandradero/sdk-go/pkg/paymentmethod"
)

func main() {
	mp.SetAccessToken("TEST-640110472259637-071923-a761f639c4eb1f0835ff7611f3248628-793910800")

	pmc := paymentmethod.NewClient()

	ch := http.Header{}
	ch.Add("X-Idempotency-Key", "some_unique_key")
	ch.Add("Some-Key", "some_value")
	opts := []mp.Option{
		mp.WithCustomHeaders(ch), // rest client will use these custom headers
	}

	res, err := pmc.List(opts...)
	if err != nil {
		panic(err)
	}

	for _, v := range res {
		fmt.Println(v)
	}
}
