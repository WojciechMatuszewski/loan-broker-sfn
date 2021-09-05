.PHONY: deploy

sam:
	sam build
	sam deploy

outputs:
	$(eval ENDPOINT=$(shell go run ./bin/outputs.go))
	@echo '{"endpoint": "${ENDPOINT}"}' > outputs.json

deploy: sam outputs
