package main

import (
	"fmt"
	"time"

	"github.com/gdeandradero/sdk-go/pkg/mp"
	"github.com/gdeandradero/sdk-go/pkg/paymentmethod"
)

func main() {
	mp.SetAccessToken("TEST-640110472259637-071923-a761f639c4eb1f0835ff7611f3248628-793910800")

	pmc := paymentmethod.NewClient()

	opts := []mp.Option{
		mp.WithTimeout(time.Second * 10), // request timeout will be 10 seconds
	}

	res, err := pmc.List(opts...)
	if err != nil {
		panic(err)
	}

	for _, v := range res {
		fmt.Println(v)
	}
}
