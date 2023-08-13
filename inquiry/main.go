package main

import (
	"os"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
)

type ContactForm struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var contactForm ContactForm 
	err := json.Unmarshal([]byte(request.Body), &contactForm)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid JSON input",
		}, nil
	}

	message := fmt.Sprintf("下記内容でお問い合わせがありました:\n\nメールアドレス:\n%s\n\n名前:\n%s\n\n件名:\n%s\n\n本文:\n%s",
		contactForm.Email, contactForm.Name, contactForm.Subject, contactForm.Message)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"), 
	})
	if err != nil {
		log.Println("Error creating AWS session:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	// Create an SNS client
	svc := sns.New(sess)

	// Replace with your SNS topic ARN
	topicArn := os.Getenv("SNS_TOPIC_ARN")

	// Publish message to the SNS topic
	_, err = svc.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(message),
	})
	if err != nil {
		log.Println("Error publishing message to SNS:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Message sent successfully",
	}, nil
}

func main() {
	lambda.Start(handler)
}
