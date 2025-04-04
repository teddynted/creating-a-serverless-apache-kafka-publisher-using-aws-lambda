package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// KafkaProducer function to send messages to Kafka
func KafkaProducer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse form data from the request body
	parsedFormData, err := url.ParseQuery(request.Body)
	if err != nil {
		log.Printf("Error parsing form data: %s", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	// Access POST parameters
	key1 := parsedFormData.Get("key1")
	key2 := parsedFormData.Get("key2")

	log.Printf("Received Key1: %s", key1)
	log.Printf("Received Key2: %s", key2)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf(`%s %s`, key1, key2),
	}, nil
}

func main() {
	lambda.Start(KafkaProducer)
}
