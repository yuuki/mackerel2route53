PKG := github.com/yuuki/mackerel2route53/src
GATEWAY_FUNC_NAME := mackerel-webhook-gateway
AUTHORIZER_FUNC_NAME := mackerel-webhook-authorizer
INSTALL_DIR := ../build

.PHONY: dep
dep:
	dep ensure

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ../build/$(GATEWAY_FUNC_NAME) -ldflags '-s -w' $(PKG)/cmd/$(GATEWAY_FUNC_NAME)
	GOOS=linux GOARCH=amd64 go build -o ../build/$(AUTHORIZER_FUNC_NAME) -ldflags '-s -w' $(PKG)/cmd/$(AUTHORIZER_FUNC_NAME)