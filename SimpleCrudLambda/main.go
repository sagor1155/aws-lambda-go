// main.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
)

const TABLE_NAME = "product-inventory"

var db *dynamodb.DynamoDB

type Product struct {
	ProductId string `json:"productId,omitempty"`
	Name      string `json:"Name,omitempty"`
	Quantity  int    `json:"Quantity,string,omitempty"`
	Brand     string `json:"Brand,omitempty"`
}

func init() {
	db = connectDynamo()
	fmt.Println("DynamoDB Initialized")
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// print context info
	fmt.Println("Lambda Invoked: ", lambdacontext.FunctionName)

	fmt.Println("------------------------------")
	fmt.Println(request)
	fmt.Println("------------------------------")

	// api path
	healthPath := "/health"
	productPath := "/product"
	productsPath := "/products"

	// get table name from environment variable
	tableName := os.Getenv("TABLE_NAME")
	fmt.Println("Table name: ", tableName)

	var response events.APIGatewayProxyResponse

	switch true {
	case request.HTTPMethod == "GET" && request.Path == healthPath:
		response = buildResponse(200, "Health Check OK")

	case request.HTTPMethod == "GET" && request.Path == productsPath:
		products, err := GetProducts()
		if err != nil {
			response = buildResponse(400, "Failed to retrieve products!")
		} else {
			data, _ := json.Marshal([]Product(products))
			response = buildResponse(200, string(data))
		}

	case request.HTTPMethod == "GET" && request.Path == productPath:
		if _, ok := request.QueryStringParameters["productId"]; ok {
			productId := request.QueryStringParameters["productId"]
			fmt.Println("product id:", productId)
			product, err := GetProduct(productId)
			if err != nil {
				response = buildResponse(404, "Failed to find requested product: "+productId)
			} else {
				msg := "ID: " + product.ProductId + "\nName: " + product.Name + "\nBrand: " + product.Brand + "\nQuantity: " + strconv.Itoa(product.Quantity)
				response = buildResponse(200, msg)
			}
		} else {
			fmt.Println("Invalid Request! productId in query parameter is missing")
			response = buildResponse(400, "Invalid Request! productId in query parameter is missing")
		}

	case request.HTTPMethod == "POST" && request.Path == productPath:
		fmt.Println(request.Body)
		product := Product{}
		json.Unmarshal([]byte(request.Body), &product)

		fmt.Println("Add Product:")
		fmt.Println("ID:  ", product.ProductId)
		fmt.Println("Name: ", product.Name)
		fmt.Println("Brand:  ", product.Brand)
		fmt.Println("Quantity:", product.Quantity)

		err := SaveProduct(product)
		if err != nil {
			response = buildResponse(400, "Failed to add product!")
		} else {
			response = buildResponse(200, "Product added successfully")
		}

	case request.HTTPMethod == "DELETE" && request.Path == productPath:
		if _, ok := request.QueryStringParameters["productId"]; ok {
			productId := request.QueryStringParameters["productId"]
			fmt.Println("product id:", productId)
			err := DeleteProduct(productId)
			if err != nil {
				response = buildResponse(400, "Failed to delete product: "+productId)
			} else {
				response = buildResponse(200, "Product deleted successfully: "+productId)
			}
		} else {
			fmt.Println("Invalid Request! productId in query parameter is missing")
			response = buildResponse(400, "Invalid Request! productId in query parameter is missing")
		}

	case request.HTTPMethod == "PATCH" && request.Path == productPath:
		fmt.Println(request.Body)
		product := Product{}
		json.Unmarshal([]byte(request.Body), &product)

		fmt.Println("Update Product:")
		fmt.Println("ID:  ", product.ProductId)
		fmt.Println("Name: ", product.Name)
		fmt.Println("Brand:  ", product.Brand)
		fmt.Println("Quantity:", product.Quantity)

		err := UpdateProduct(product)
		if err != nil {
			response = buildResponse(400, "Failed to update product!")
		} else {
			response = buildResponse(200, "Product: "+string(product.ProductId)+" updated successfully")
		}

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

func UpdateProduct(product Product) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(TABLE_NAME),

		Key: map[string]*dynamodb.AttributeValue{ // which product
			"productId": {
				S: aws.String(product.ProductId),
			},
		},

		UpdateExpression:// update expression
		aws.String("set #nm = :n, #br = :b, #qt = :q"),

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{ // update values
			":n": {
				S: aws.String(product.Name),
			},
			":b": {
				S: aws.String(product.Brand),
			},
			":q": {
				S: aws.String(strconv.Itoa(product.Quantity)),
			},
		},

		ExpressionAttributeNames: map[string]*string{ // update values
			"#nm": aws.String("Name"),
			"#br": aws.String("Brand"),
			"#qt": aws.String("Quantity"),
		},

		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err := db.UpdateItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	return err
}

func DeleteProduct(productId string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"productId": {
				S: aws.String(productId),
			},
		},
		TableName: aws.String(TABLE_NAME),
	}

	_, err := db.DeleteItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	return err
}

func GetProducts() ([]Product, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(TABLE_NAME),
	}

	result, err := db.Scan(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	products := []Product{}

	for _, item := range result.Items {
		product := Product{}
		err = dynamodbattribute.UnmarshalMap(item, &product)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		} else {
			products = append(products, product)
		}
	}

	for _, product := range products {
		fmt.Println("-----------------------------")
		fmt.Println("ID:  ", product.ProductId)
		fmt.Println("Name: ", product.Name)
		fmt.Println("Brand:  ", product.Brand)
		fmt.Println("Quantity:", product.Quantity)
		fmt.Println("-----------------------------")
	}

	return products, nil
}

func GetProduct(productId string) (Product, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"productId": {
				S: aws.String(productId),
			},
		},
		TableName: aws.String(TABLE_NAME),
	}

	result, err := db.GetItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	product := Product{}

	if result.Item == nil {
		msg := "Could not find product: " + productId
		return product, errors.New(msg)
	}

	// fmt.Println(result.Item)
	err = dynamodbattribute.UnmarshalMap(result.Item, &product)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	fmt.Println("Found Product:")
	fmt.Println("ID:  ", product.ProductId)
	fmt.Println("Name: ", product.Name)
	fmt.Println("Brand:  ", product.Brand)
	fmt.Println("Quantity:", product.Quantity)

	return product, nil
}

func SaveProduct(product Product) error {
	input := &dynamodb.PutItemInput{
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
	}

	_, err := db.PutItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	return err
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
