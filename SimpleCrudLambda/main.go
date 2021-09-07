// main.go
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
)

const TABLE_NAME = "product-inventory"
var db *dynamodb.DynamoDB

type Product struct {
	ProductId string
	Name string
	Quantity int
	Brand string
}

func init(){
	db = connectDynamo()
	fmt.Println("DynamoDB Initialized")
}

func connectDynamo() (db *dynamodb.DynamoDB) {
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println("Failed to initialize aws session!")
		os.Exit(1)
	}
	return dynamodb.New(sess)
}

func SaveProduct(product Product) {
	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"productId": {
				S: aws.String(product.ProductId),
			},
			"Name": {
				S: aws.String(product.Name),
			},
			"Quantity": {
				S: aws.String(strconv.Itoa(product.Quantity)),
			},
			"Brand": {
				S: aws.String(product.Brand),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
}

func buildResponse(status int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	// product := Product{
	// 	ProductId: "P1",
	// 	Name: "car",
	// 	Quantity: 2,
	// 	Brand: "volkswagon",
	// }
	// SaveProduct(product)

	// print context info
	lc, _ := lambdacontext.FromContext(ctx)
	fmt.Println("Context:", lc.Identity.CognitoIdentityPoolID)

	// api path
	healthPath := "/health"
	productPath := "/product"
	productsPath := "/products"

	// get table name from environment variable
	// tableName := os.Getenv("TABLE_NAME")
	// fmt.Println("Table name: ", tableName)

	var response events.APIGatewayProxyResponse

	switch true {
	case request.HTTPMethod == "GET" && request.Path == healthPath:
		fmt.Println("GET: healthpath")
		response = buildResponse(200, "GET: healthpath")

	case request.HTTPMethod == "GET" && request.Path == productsPath:
		fmt.Println("GET: productsPath")
		response = buildResponse(200, "GET: productsPath")

	case request.HTTPMethod == "GET" && request.Path == productPath:
		fmt.Println("GET: productPath")
		response = buildResponse(200, "GET: productPath")

	case request.HTTPMethod == "POST" && request.Path == productPath:
		fmt.Println("POST: productPath")
		response = buildResponse(200, "POST: productPath")

	case request.HTTPMethod == "DELETE" && request.Path == productPath:
		fmt.Println("DELETE: productPath")
		response = buildResponse(200, "DELETE: productPath")

	case request.HTTPMethod == "PATCH" && request.Path == productPath:
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
