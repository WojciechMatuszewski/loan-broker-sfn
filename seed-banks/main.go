package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"golang.org/x/net/context"
)

func main() {
	lambda.Start(cfn.LambdaWrap(handler))
}

type Bank struct {
	MaxLoanAmount  int    `json:"maxLoanAmount,omitempty"`
	MinCreditScore int    `json:"minCreditScore,omitempty"`
	BaseRate       int    `json:"baseRate,omitempty"`
	ID             string `json:"id,omitempty"`
}

type BankItem struct {
	PK             string `dynamodbav:"pk"`
	MaxLoanAmount  int    `dynamodbav:"maxLoanAmount,omitempty"`
	MinCreditScore int    `dynamodbav:"minCreditScore,omitempty"`
	BaseRate       int    `dynamodbav:"baseRate,omitempty"`
}

var banks = []Bank{
	{ID: "Universal", BaseRate: 4, MaxLoanAmount: 700000, MinCreditScore: 500},
	{ID: "PawnShop", BaseRate: 5, MaxLoanAmount: 500000, MinCreditScore: 400},
	{ID: "Premium", BaseRate: 3, MaxLoanAmount: 900000, MinCreditScore: 600},
}

func handler(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	fmt.Println(event)

	if err != nil {
		return physicalResourceID, nil, err
	}

	tableName, found := event.ResourceProperties["TableName"].(string)
	if !found {
		return physicalResourceID, nil, errors.New("TableName parameter not found")
	}

	if event.RequestType == cfn.RequestUpdate {
		return
	}

	if event.RequestType == cfn.RequestDelete {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	db := dynamodb.NewFromConfig(cfg)

	if event.RequestType == cfn.RequestCreate {
		for _, bank := range banks {
			fmt.Println("Creating bank", bank.ID)

			item := BankItem{
				PK:             fmt.Sprintf("BANK#%v", bank.ID),
				MaxLoanAmount:  bank.MaxLoanAmount,
				MinCreditScore: bank.MinCreditScore,
				BaseRate:       bank.BaseRate,
			}

			avs, err := attributevalue.MarshalMap(item)
			if err != nil {
				return physicalResourceID, nil, err
			}

			expr, err := expression.NewBuilder().WithCondition(expression.AttributeNotExists(expression.Name("pk"))).Build()
			if err != nil {
				return physicalResourceID, nil, err
			}

			_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
				Item:                      avs,
				TableName:                 aws.String(tableName),
				ConditionExpression:       expr.Condition(),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
			})
			if err != nil {
				return physicalResourceID, nil, err
			}
		}
	}

	return
}
