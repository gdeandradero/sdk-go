package main

import (
	"fmt"
	"time"

	"github.com/gdeandradero/sdk-go/pkg/config"
	"github.com/gdeandradero/sdk-go/pkg/payment"
)

func main() {
	config.New("TEST-640110472259637-071923-a761f639c4eb1f0835ff7611f3248628-793910800")

	id := createPayment()
	time.Sleep(2 * time.Second)
	searchPayment()
	getPayment(id)
}

var externalReference = "312931sadaddasddaasdadsaddsa29asdasdasdsaasdadsadsa3919321"

func createPayment() int64 {
	pc := payment.NewClient()

	request := payment.Request{
		TransactionAmount: 1.5,
		PaymentMethodID:   "pix",
		Description:       "meu pagamento",
		ExternalReference: externalReference,
		Payer: &payment.PayerRequest{
			Email: "fhashfadsuhfdafasdfasfashfda@testuser.com",
		},
	}
	res, err := pc.Create(request)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	fmt.Println(res.ID)
	return res.ID
}

func searchPayment() {
	pc := payment.NewClient()

	request := payment.Filters{
		Sort:              "date_created",
		Criteria:          "asc",
		ExternalReference: externalReference,
		Range:             "date_created",
		BeginDate:         "2023-01-01T00:00:00Z",
		EndDate:           "2024-01-01T00:00:00Z",
	}
	res, err := pc.Search(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.Results[0].ID)
}

func getPayment(id int64) {
	pc := payment.NewClient()

	res, err := pc.Get(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.ID)
}
