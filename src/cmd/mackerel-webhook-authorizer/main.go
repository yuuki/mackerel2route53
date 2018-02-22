package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func generatePolicy(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "mackerel-webhook-authorizer",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action: []string{"execute-api:Invoke"},
					Effect: "Allow",
					Condition: events.TestCondition{
						IpAddress: map[string][]string{
							"aws:SourceIp": []string{
								"52.193.111.118", "52.196.125.133", "13.113.213.40", "52.197.186.229", "52.198.79.40", "13.114.12.29", "13.113.240.89", "52.68.245.9", "13.112.142.176",
							},
						},
					},
					Resource: []string{event.MethodArn},
				},
			},
		},
	}, nil
}

func main() {
	lambda.Start(generatePolicy)
}
