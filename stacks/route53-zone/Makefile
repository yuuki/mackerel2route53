PARAMETERS := '$(shell printf "cat <<++EOS\n`cat ./parameters.json`\n++EOS\n" | bash | tr -d "\n")'
TEMPLATE := "file://$(PWD)/template.yml"
STACK_NAME := mackerel2route53-zones

.PHONY: create
create:
	aws cloudformation create-stack --stack-name $(STACK_NAME) --parameters $(PARAMETERS) --template-body $(TEMPLATE) --capabilities CAPABILITY_NAMED_IAM
	aws cloudformation wait stack-create-complete --stack-name $(STACK_NAME)

.PHONY: update
update:
	aws cloudformation update-stack --stack-name $(STACK_NAME) --parameters $(PARAMETERS) --template-body $(TEMPLATE) --capabilities CAPABILITY_NAMED_IAM
	cloudformation wait stack-update-complete --stack-name $(STACK_NAME)

.PHONY: destroy
destroy:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)