package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"

	storage "github.com/shjp/shjp-storage"
	"github.com/shjp/shjp-storage/s3"
)

// Config is the config parameters
type Config struct {
	S3Region          string
	S3Bucket          string
	S3AccessKeyID     string
	S3SecretAccessKey string
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
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

	err = svc.Upload("foofoo", nil)
	if err != nil {
		return formatResponse(http.StatusInternalServerError, err.Error()), err
	}

	return formatResponse(http.StatusOK, "foo"), nil
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
