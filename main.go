package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// MackerelWebhookRequest represents a webhook request for Mackerel.
// https://mackerel.io/ja/docs/entry/howto/alerts/webhook
type MackerelWebhookRequest struct {
	OrgName string               `json:"orgName"`
	Event   string               `json:"event"`
	Host    *MackerelWebhookHost `json:"host"`
}

// MackerelWebhookHost contains identity information for the Mackerel host.
type MackerelWebhookHost struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	URL       string                 `json:"url"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	Memo      string                 `json:"memo"`
	IsRetired bool                   `json:"isRetired"`
	Roles     []*MackerelWebhookRole `json:"roles"`
}

// MackerelWebhookRole contains identity information for the Mackerel host's role.
type MackerelWebhookRole struct {
	Fullname    string `json:"fullname"`
	ServiceName string `json:"serviceName"`
	ServiceURL  string `json:"serviceUrl"`
	RoleName    string `json:"roleName"`
	RoleURL     string `json:"roleUrl"`
}

// Response contains response's message.
type Response struct {
	Message string
}

var (
	svc *route53.Route53
)

func init() {
	svc = route53.New(session.New())
}

func mackerelWebhookHandler(req MackerelWebhookRequest) (Response, error) {
	return Response{Message: "success"}, nil
}
func main() {
	lambda.Start(mackerelWebhookHandler)
}
