# Makefile for project
.PHONY: all

all: help

## Help
help: ## 
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    %-20s%s\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  \n%s\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

## Test
test: ## 
	make -i build
	make -i upload-s3
	make -i update-lambda

build: ## 
	cd lambda-api; GOARCH=amd64 GOOS=linux go build .

upload-s3: ## 
	cd lambda-api; zip gs-aws-lambda-demo.zip gs-aws-lambda-demo
	cd lambda-api; aws s3 cp gs-aws-lambda-demo.zip s3://gs-aws-lambda-demo

update-lambda: ## 
	aws lambda update-function-code \
    --function-name gs-aws-lambda-demo \
    --s3-bucket gs-aws-lambda-demo \
    --s3-key gs-aws-lambda-demo.zip
