module loan-broker

go 1.16

require (
	github.com/aws/aws-lambda-go v1.26.0
	github.com/aws/aws-sdk-go-v2 v1.9.0
	github.com/aws/aws-sdk-go-v2/config v1.8.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.2.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.2.2
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.10.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.5.0
	github.com/pelletier/go-toml v1.9.3
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
