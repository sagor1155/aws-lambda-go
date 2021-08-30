// main.go
package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s]: %s Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, record.EventName, s3.Bucket.Name, s3.Object.Key)

	}

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	sourceName := "s3bucket-source"
	destinationName := "s3bucket-destination"
	objectKey := s3Event.Records[0].S3.Object.Key
	copySource := sourceName + "/" + objectKey

	// Copy the item
	_, err = svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(destinationName),
		CopySource: aws.String(url.PathEscape(copySource)), Key: aws.String(objectKey)})

	if err != nil {
		fmt.Printf("Unable to copy item from bucket %q to bucket %q, %v", sourceName, destinationName, err)
	} else {
		fmt.Println("S3 Bucket Copy Successfull")
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
