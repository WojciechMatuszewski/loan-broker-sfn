package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

type Event struct {
	Region                     string                 `json:"region,omitempty"`
	AccountId                  string                 `json:"account_id,omitempty"`
	Fragment                   map[string]interface{} `json:"fragment,omitempty"`
	TransformId                string                 `json:"transformId,omitempty"`
	Params                     interface{}            `json:"params,omitempty"`
	RequestId                  string                 `json:"requestId,omitempty"`
	TemplateParameterVariables interface{}            `json:"templateParameterValues,omitempty"`
}

type Response struct {
	RequestId string                 `json:"requestId,omitempty"`
	Status    string                 `json:"status,omitempty"`
	Fragment  map[string]interface{} `json:"fragment,omitempty"`
}

func handler(ctx context.Context, event Event) (Response, error) {
	fmt.Println("------------------------")

	buf, err := json.Marshal(event.Fragment)
	if err != nil {
		panic(err)
	}

	fmt.Println("Old fragment", string(buf))

	timestamp := strconv.Itoa(int(time.Now().Unix()))
	fmt.Println("Timestamp", timestamp)

	rawNewFragment := strings.ReplaceAll(string(buf), "$TIMESTAMP$", timestamp)
	fmt.Println("New changed fragment", rawNewFragment)

	var newFragment map[string]interface{}
	err = json.Unmarshal([]byte(rawNewFragment), &newFragment)
	if err != nil {
		panic(err)
	}

	fmt.Println("------------------------")

	return Response{
		RequestId: event.RequestId,
		Status:    "success",
		Fragment:  newFragment,
	}, nil
}
