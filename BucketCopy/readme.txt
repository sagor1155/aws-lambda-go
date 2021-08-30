install dependencies:
    go get -u github.com/aws/aws-sdk-go
    go get github.com/aws/aws-lambda-go/lambda
    go get github.com/aws/aws-lambda-go/events
    go get github.com/aws/aws-lambda-go/lambdacontext

go module init:
    go mod init your-module-name
    go mod tidy

build: 
    GOOS=linux GOARCH=amd64 go build -o main main.go

make zip: 
    zip project-name.zip executable-name
