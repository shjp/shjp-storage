package main

import (
	"encoding/json"
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
	errReadingConfig    = errors.New("can't read config")
	errMethodNotAllowed = errors.New("only POST method allowed")
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
	reqBlob, err := json.Marshal(request)
	if err != nil {
		log.Println("Marshalling request failed:", err)
	}
	log.Println("Request object ---------------------------------------------------")
	log.Println(string(reqBlob))
	log.Println("------------------------------------------------------------------")

	if request.HTTPMethod == http.MethodOptions {
		return formatPreflightResponse(), nil
	}
	if request.HTTPMethod != http.MethodPost {
		return formatResponse(http.StatusMethodNotAllowed, "Disallowed HTTP verb"), errMethodNotAllowed
	}

	var config Config
	err = envconfig.Process("aws", &config)
	if err != nil {
		log.Printf("Error reading config: %s\n", err)
		return formatResponse(http.StatusInternalServerError, "Server error"), errReadingConfig
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

	return formatResponse(http.StatusOK, `{"url":`+url+`}`), nil
}

func formatResponse(statusCode int, body string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST,OPTIONS",
		},
		Body: body,
	}
}

func formatPreflightResponse() *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
		Body: "{}",
	}
}
