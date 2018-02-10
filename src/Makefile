FUNC_NAME := mackerel-webhook-gateway

.PHONY: dep
dep:
	dep ensure

.PHONY: build-function
build-function:
	GOOS=linux GOARCH=amd64 go build -o ../build/$(FUNC_NAME) -ldflags '-s -w' ./...

