package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

func Handler(ctx context.Context, req Request) (Response, error) {
	fmt.Println("Path:", req.Path)
	fmt.Println("Method:", req.HTTPMethod)

	body := map[string]string{
		"message": "Lambda executed successfully",
		"status":  "success",
	}

	b, _ := json.Marshal(body)

	resp := Response{
		StatusCode: 200,
		Body:       string(b),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
