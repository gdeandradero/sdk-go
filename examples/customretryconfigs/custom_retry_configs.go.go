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
		mp.WithMaxRetries(3),                      // default retry client will retry 3 times
		mp.WithMaxBackoff(time.Second * 30),       // default retry client will wait 30 seconds til the next request
		mp.WithRetryDelay(time.Millisecond * 500), // default retry client initial delay will have 0.5 seconds
	}

	res, err := pmc.List(opts...)
	if err != nil {
		panic(err)
	}

	for _, v := range res {
		fmt.Println(v)
	}
}
