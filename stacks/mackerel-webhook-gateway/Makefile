PROJECT := mackerel2route53
FUNC_NAME := mackerel-webhook-gateway
GATEWAY_ID := $(shell aws apigateway get-rest-apis | jq -rM '.items[] | select(.name == "$(PROJECT)") | .id')
GATEWAY_URL := https://$(GATEWAY_ID).execute-api.ap-northeast-1.amazonaws.com/Prod

.PHONY: build
build:
	(cd ../../src; make build)

.PHONY: deploy
deploy:
	@../../script/envfile ../../config.env make _deploy

.PHONY: _deploy
_deploy: build
	aws cloudformation package \
		--template-file template.yml \
		--s3-bucket $(MACKEREL2ROUTE53_FUNCTION_S3_BUCKET) \
		--s3-prefix $(PROJECT) \
		--output-template-file .template.yml
	aws cloudformation deploy \
		--template-file .template.yml \
		--parameter-overrides 'MackerelAPIKey=$(MACKEREL2ROUTE53_MACKEREL_API_KEY)' \
		--stack-name $(PROJECT) \
		--capabilities CAPABILITY_IAM

.PHONY: destroy
destroy:
	aws cloudformation delete-stack --stack-name $(PROJECT)

.PHONY: logs
logs:
	awslogs get "/aws/lambda/$(FUNC_NAME)" --watch

.PHONY: getway-url
gateway-url:
	@echo $(GATEWAY_URL)

.PHONY: local-invoke
local-invoke: build
	sam local invoke 'App'