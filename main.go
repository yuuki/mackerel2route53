package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

// MackerelWebhookRequest represents a webhook request for Mackerel.
// https://mackerel.io/ja/docs/entry/howto/alerts/webhook
type MackerelWebhookRequest struct {
	OrgName     string `json:"orgName"`
	Event       string `json:"event"`
	WebhookHost string `json:"host"`
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

type Response struct {
	Message string
}

func mackerelWebhookHandler(req MackerelWebhookRequest) (Response, error) {
	return Response{Message: "success"}, nil
}
func main() {
	lambda.Start(mackerelWebhookHandler)
}
