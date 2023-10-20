package main

import (
	"fmt"

	"github.com/gdeandradero/sdk-go/pkg/mp"
	"github.com/gdeandradero/sdk-go/pkg/paymentmethod"
)

func main() {
	rc := mp.NewRestClient("TEST-640110472259637-071923-a761f639c4eb1f0835ff7611f3248628-793910800")

	pmc := paymentmethod.NewClient(rc)
	res, err := pmc.List()
	if err != nil {
		panic(err)
	}

	for _, v := range res {
		fmt.Println(v)
	}
}
