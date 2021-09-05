package main

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

type Event struct {
	BankInfo struct {
		MinCreditScore string `json:"minCreditScore,omitempty"`
		BaseRate       string `json:"baseRate,omitempty"`
		MaxLoanAmount  string `json:"maxLoanAmount,omitempty"`
	} `json:"bankInfo,omitempty"`
	Amount   float64 `json:"amount,omitempty"`
	Term     float64 `json:"term,omitempty"`
	BankName string  `json:"bankName"`
	Credit   struct {
		Score   float64 `json:"score,omitempty"`
		History float64 `json:"history,omitempty"`
	} `json:"credit,omitempty"`
}

type Output struct {
	Rate     float64 `json:"rate"`
	BankName string  `json:"bankName"`
}

func handler(ctx context.Context, event Event) (*Output, error) {
	maxLoanAmount, err := strconv.ParseFloat(event.BankInfo.MaxLoanAmount, 64)
	if err != nil {
		return nil, err
	}

	minCreditScore, err := strconv.ParseFloat(event.BankInfo.MinCreditScore, 64)
	if err != nil {
		return nil, err
	}

	if event.Amount > maxLoanAmount || event.Credit.Score < minCreditScore {
		return nil, nil
	}

	baseRate, err := strconv.ParseFloat(event.BankInfo.BaseRate, 64)
	if err != nil {
		return nil, err
	}

	rate := baseRate*rand.Float64() + ((1000 - event.Credit.Score) / 100.0)
	return &Output{Rate: rate, BankName: event.BankName}, nil
}
