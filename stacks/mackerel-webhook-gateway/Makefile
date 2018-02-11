PROJECT := mackerel2route53
FUNC_NAME := mackerel-webhook-gateway
GATEWAY_ID := $(shell aws apigateway get-rest-apis | jq -rM '.items[] | select(.name == "$(PROJECT)") | .id')
GATEWAY_URL := https://$(GATEWAY_ID).execute-api.ap-northeast-1.amazonaws.com/Prod

.PHONY: deploy
deploy:
	(cd ../../src; make build)
	aws cloudformation package \
		--template-file template.yml \
		--s3-bucket yuuki-lambda-packages \
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
local-invoke:
	../../script/envfile ../../config.env sam local invoke 'App'