package main_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSFNContainer(t *testing.T) {
	ctx := context.Background()
	port, err := nat.NewPort("tcp", "8083")
	if err != nil {
		t.Fatalf(err.Error())
	}

	req := testcontainers.ContainerRequest{
		Image:        "amazon/aws-stepfunctions-local",
		ExposedPorts: []string{"8083/tcp"},
		WaitingFor:   wait.ForListeningPort(port),
	}

	sfnContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer sfnContainer.Terminate(ctx)

	sfnBuf, err := os.ReadFile("/Users/wojciechmatuszewsk-personal/Desktop/currently-learning/loan-broker-go/broker-machine.asl.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	sfnDef := strings.ReplaceAll(string(sfnBuf), "${DataTableName}", "test123")

	endpoint, err := sfnContainer.Endpoint(ctx, "http")
	if err != nil {
		t.Fatal(err)
	}

	sfnClient := sfn.New(sfn.Options{EndpointResolver: sfn.EndpointResolverFromURL(endpoint)})
	cOut, err := sfnClient.CreateStateMachine(ctx, &sfn.CreateStateMachineInput{
		Definition: aws.String(sfnDef),
		Name:       aws.String("testMachine"),
		RoleArn:    aws.String("arn:aws:iam::012345678901:role/DummyRole"),
	})
	if err != nil {
		t.Fatal(err)
	}

	sOut, err := sfnClient.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: cOut.StateMachineArn,
		Input:           aws.String("{\"foo\": \"bar\"}"),
		Name:            aws.String("wow"),
	})
	if err != nil {
		t.Fatal(err)
	}

	for {
		dOut, err := sfnClient.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: sOut.ExecutionArn,
		})
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(dOut.Output, dOut.Status)

		if dOut.Status != types.ExecutionStatusRunning {
			break
		}

		fmt.Println("sleeping")
		time.Sleep(time.Duration(time.Second * 1))
	}

	paginator := sfn.NewGetExecutionHistoryPaginator(sfnClient, &sfn.GetExecutionHistoryInput{ExecutionArn: sOut.ExecutionArn})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for _, event := range output.Events {
			spew.Dump(event)
		}
	}

}
