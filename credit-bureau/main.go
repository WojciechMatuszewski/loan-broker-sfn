package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

type Event struct {
	SSN    string `json:"SSN,omitempty"`
	Amount int    `json:"amount,omitempty"`
	Term   int    `json:"term,omitempty"`
}

type Output struct {
	Score   int `json:"score,omitempty"`
	History int `json:"history,omitempty"`
}

func handler(ctx context.Context, event Event) (Output, error) {
	fmt.Println(event)

	const minScore = 300
	const maxScore = 900

	return Output{
		Score:   rand.Intn((maxScore - minScore) + minScore),
		History: rand.Intn((30 - 1) + 1),
	}, nil
}
