package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	mkr "github.com/mackerelio/mackerel-client-go"
)

// MackerelWebhookRequest represents a webhook request for Mackerel.
// https://mackerel.io/ja/docs/entry/howto/alerts/webhook
type MackerelWebhookRequest struct {
	OrgName    string               `json:"orgName"`
	Event      string               `json:"event"`
	Host       *MackerelWebhookHost `json:"host"`
	FromStatus string               `json:"fromStatus"`
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

const (
	dnsRecordTTL = 300
)

var (
	mackerelAPIKey string
	zoneID         string
	svc            *route53.Route53
)

func init() {
	if mackerelAPIKey = os.Getenv("MACKEREL2ROUTE53_MACKEREL_API_KEY"); mackerelAPIKey == "" {
		log.Println("MACKEREL2ROUTE53_MACKEREL_API_KEY is empty")
		os.Exit(1)
	}
	if zoneID = os.Getenv("MACKEREL2ROUTE53_ZONE_ID"); zoneID == "" {
		log.Println("MACKEREL2ROUTE53_ZONE_ID is empty")
		os.Exit(1)
	}
	svc = route53.New(session.New())
}

func response(code int, msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       fmt.Sprintf("{\"message\":\"%s\"}", msg),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func findIPAddressFromMackerel(hostID string) (string, error) {
	client := mkr.NewClient(mackerelAPIKey)
	host, err := client.FindHost(hostID)
	if err != nil {
		return "", err
	}
	if len(host.Interfaces) == 0 {
		return "", fmt.Errorf("Not found interfaces on Mackerel. host: %v", hostID)
	}
	return host.Interfaces[0].IPAddress, nil // Adapt the first interface.
}

func createRecord(host *MackerelWebhookHost) error {
	log.Printf("createRecord: %v\n", host)

	ipaddr, err := findIPAddressFromMackerel(host.ID)
	if err != nil {
		return err
	}

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("CREATE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(host.Name),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ipaddr),
							},
						},
						TTL: aws.Int64(dnsRecordTTL),
					},
				},
			},
		},
		HostedZoneId: aws.String(zoneID),
	}
	_, err = svc.ChangeResourceRecordSets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println(aerr)
			switch aerr.Code() {
			case route53.ErrCodePriorRequestNotComplete:
				//TODO wait&retry?
			}
		}
		return err
	}

	return nil
}

func updateRecord(host *MackerelWebhookHost, fromStatus string) error {
	log.Printf("updateRecord: %+v, %v\n", host, fromStatus)

	return nil
}

func deleteRecord(host *MackerelWebhookHost) error {
	log.Printf("deleteRecord: %+v\n", host)

	ipaddr, err := findIPAddressFromMackerel(host.ID)
	if err != nil {
		return err
	}

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("DELETE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(host.Name),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ipaddr),
							},
						},
						TTL: aws.Int64(dnsRecordTTL),
					},
				},
			},
		},
		HostedZoneId: aws.String(zoneID),
	}
	if _, err := svc.ChangeResourceRecordSets(input); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println(aerr)
			switch aerr.Code() {
			case route53.ErrCodePriorRequestNotComplete:
				//TODO wait&retry?
			}
		}
		return err
	}

	return nil
}

func mackerelWebhookHandler(ctx context.Context, gwReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req MackerelWebhookRequest
	if err := json.Unmarshal([]byte(gwReq.Body), &req); err != nil {
		return response(400, "json decode error"), err
	}

	log.Printf("%+v\n", req)

	switch req.Event {
	case "hostRegister":
		if err := createRecord(req.Host); err != nil {
			return response(500, "Route53 record creation error"), err
		}
	case "hostStatus":
		if err := updateRecord(req.Host, req.FromStatus); err != nil {
			return response(500, "Route53 record update error"), err
		}
	case "hostRetire":
		if err := deleteRecord(req.Host); err != nil {
			return response(500, "Route53 record delete error"), err
		}
	default:
		err := fmt.Errorf("invalid webhook request event: %s", req.Event)
		return response(400, err.Error()), err
	}
	return response(200, "success"), nil
}
func main() {
	lambda.Start(mackerelWebhookHandler)
}
