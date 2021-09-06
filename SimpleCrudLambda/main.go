// main.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func buildResponse(status int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	_, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		fmt.Println("Failed to initialize aws session!")
		os.Exit(1)
	}

	// print context info
	lc, _ := lambdacontext.FromContext(ctx)
	fmt.Println(lc.Identity.CognitoIdentityPoolID)

	// api path
	healthPath := "/health"
	productPath := "/product"
	productsPath := "/products"

	// get table name from environment variable
	tableName := os.Getenv("TABLE_NAME")
	fmt.Println("Table name: ", tableName)

	// create DynamoDB client
	// db = dynamodb.New(sess)

	var response events.APIGatewayProxyResponse

	switch true {
	case event.HTTPMethod == "GET" && event.Path == healthPath:
		fmt.Println("GET: healthpath")
		response = buildResponse(200, "GET: healthpath")

	case event.HTTPMethod == "GET" && event.Path == productsPath:
		fmt.Println("GET: productsPath")
		response = buildResponse(200, "GET: productsPath")

	case event.HTTPMethod == "GET" && event.Path == productPath:
		fmt.Println("GET: productPath")
		response = buildResponse(200, "GET: productPath")

	case event.HTTPMethod == "POST" && event.Path == productPath:
		fmt.Println("POST: productPath")
		response = buildResponse(200, "POST: productPath")

	case event.HTTPMethod == "DELETE" && event.Path == productPath:
		fmt.Println("DELETE: productPath")
		response = buildResponse(200, "DELETE: productPath")

	case event.HTTPMethod == "PATCH" && event.Path == productPath:
		fmt.Println("PATCH: productPath")
		response = buildResponse(200, "PATCH: productPath")

	default:
		fmt.Println("Invalid Request!")
		response = buildResponse(404, "Bad Request")
	}

	return response, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}

/*
// APIGatewayProxyRequest contains data coming from the API Gateway proxy
type APIGatewayProxyRequest struct {
	Resource                        string                        `json:"resource"` // The resource path defined in API Gateway
	Path                            string                        `json:"path"`     // The url path for the caller
	HTTPMethod                      string                        `json:"httpMethod"`
	Headers                         map[string]string             `json:"headers"`
	MultiValueHeaders               map[string][]string           `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string             `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string           `json:"multiValueQueryStringParameters"`
	PathParameters                  map[string]string             `json:"pathParameters"`
	StageVariables                  map[string]string             `json:"stageVariables"`
	RequestContext                  APIGatewayProxyRequestContext `json:"requestContext"`
	Body                            string                        `json:"body"`
	IsBase64Encoded                 bool                          `json:"isBase64Encoded,omitempty"`
}

// APIGatewayProxyResponse configures the response to be returned by API Gateway for the request
type APIGatewayProxyResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
}
*/
