package main

import (
	"fmt"
	"log"

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
	svc *route53.Route53
)

func init() {
	svc = route53.New(session.New())
}

func findIPAddressFromMackerel(hostID string) (string, error) {
	client := mkr.NewClient("API-KEY")
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
		HostedZoneId: aws.String("Z3M3LMPEXAMPLE"), // FIXME
	}
	_, err = svc.ChangeResourceRecordSets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
			case route53.ErrCodeNoSuchHealthCheck:
				fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
			case route53.ErrCodeInvalidChangeBatch:
				fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
			case route53.ErrCodeInvalidInput:
				fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
			case route53.ErrCodePriorRequestNotComplete:
				fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return nil
}

func updateRecord(host *MackerelWebhookHost, fromStatus string) error {
	return nil
}

func deleteRecord(host *MackerelWebhookHost) error {
	return nil
}

func mackerelWebhookHandler(req MackerelWebhookRequest) (Response, error) {
	switch req.Event {
	case "hostRegister":
		if err := createRecord(req.Host); err != nil {
			return Response{Message: "failed to create Route53 record"}, err
		}
	case "hostStatus":
		if err := updateRecord(req.Host, req.FromStatus); err != nil {
			return Response{Message: "failed to update Route53 record"}, err
		}
	case "hostRetire":
		if err := deleteRecord(req.Host); err != nil {
			return Response{Message: "failed to delete Route53 record"}, err
		}
	}
	return Response{Message: "success"}, nil
}
func main() {
	lambda.Start(mackerelWebhookHandler)
}
