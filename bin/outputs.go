package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/pelletier/go-toml"
)

func main() {
	base, err := filepath.Abs("./")
	if err != nil {
		fmt.Println("err")
		panic(err)
	}

	var walker func(basePath string, level int) (string, error)
	walker = func(basePath string, level int) (string, error) {
		var absPath string

		err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				return nil
			}

			filePath := fmt.Sprintf("%v/samconfig.toml", path)
			_, err = os.Stat(filePath)

			if err == nil {
				absPath = filePath
			}

			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return "", err
		}

		if absPath != "" {
			return absPath, nil
		}

		if level == 2 {
			return "", errors.New("not found")
		}

		if level < 2 {
			return walker(filepath.Join(basePath, "../"), level+1)
		}

		return absPath, nil
	}

	absPath, err := walker(base, 0)
	if err != nil {
		panic(err)
	}

	t, err := toml.LoadFile(absPath)
	if err != nil {
		panic(err)
	}

	stackName, found := t.Get("default.deploy.parameters.stack_name").(string)
	if !found {
		panic(errors.New("stack_name not found"))
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	cfnClient := cloudformation.NewFromConfig(cfg)

	out, err := cfnClient.DescribeStacks(context.Background(), &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		panic(err)
	}

	if len(out.Stacks) == 0 {
		panic(errors.New("stack not found"))
	}

	outputs := out.Stacks[0].Outputs

	var endpointURL string
	for _, output := range outputs {
		if *output.OutputKey == "BrokerAPIEntryMethod" {
			endpointURL = *output.OutputValue
		}
	}
	if endpointURL == "" {
		panic(errors.New("not found"))
	}

	fmt.Println(endpointURL)
}
