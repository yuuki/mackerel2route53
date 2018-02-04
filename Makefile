.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/lambda-mackerel2route53 -ldflags '-s -w' ./...

.PHONY: deploy
deploy: build-linux
	aws cloudformation package \
		--template-file template.yml \
		--s3-bucket yuuki-lambda-packages \
		--s3-prefix lambda-mackerel2route53 \
		--output-template-file .template.yml
	aws cloudformation deploy \
		--template-file .template.yml \
		--stack-name lambda-mackerel2route53 \
		--capabilities CAPABILITY_IAM