package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"

	storage "github.com/shjp/shjp-storage"
	"github.com/shjp/shjp-storage/s3"
)

var (
	errMethodNotAllowed = errors.New("Only POST method allowed")
)

// Config is the config parameters
type Config struct {
	S3Region string
	S3Bucket string
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == http.MethodOptions {
		return formatResponse(http.StatusOK, "OK"), nil
	}
	if request.HTTPMethod != http.MethodPost {
		return formatResponse(http.StatusMethodNotAllowed, "Disallowed HTTP verb"), errMethodNotAllowed
	}

	var config Config
	err := envconfig.Process("aws", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Config: %#v\n", config)

	// authToken, ok := request.Headers["auth-token"]
	// if !ok {
	// 	log.Println("Auth token not found")
	// }

	log.Println("body:", request.Body)
	svc := storage.NewService(s3.Client{
		Region: config.S3Region,
		Bucket: config.S3Bucket,
	})

	url, err := svc.Upload("foofoo", "barbar", nil)
	if err != nil {
		return formatResponse(http.StatusInternalServerError, err.Error()), err
	}

	return formatResponse(http.StatusOK, url), nil
}

func formatResponse(statusCode int, body string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "X-Requested-With,Content-Type,Authorization,Auth-Token",
			"Access-Control-Allow-Methods": "GET,PUT,POST,DELETE,OPTIONS,PING",
		},
		Body: body,
	}
}
