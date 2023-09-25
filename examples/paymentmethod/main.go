package main

import (
	"fmt"

	"github.com/gdeandradero/sdk-go/pkg/config"
	"github.com/gdeandradero/sdk-go/pkg/paymentmethod"
)

func main() {
	config.New("TEST-640110472259637-071923-a761f639c4eb1f0835ff7611f3248628-793910800")

	pmc := paymentmethod.NewClient()
	x, err := pmc.List()
	if err != nil {
		panic(err)
	}

	for _, v := range x {
		fmt.Println(v)
	}
}
