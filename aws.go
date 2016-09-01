package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func testSQS() (err error) {
	creds := credentials.NewStaticCredentials("AKIAJMZJ3X6MHMXRF7QQ", "0hincCnlm2XvpdpSD+LBs6NSwfF0250pEnEyYJ49", "")
	_, err = creds.Get()
	if err != nil {
		log.Fatal(err)
	}

	cfg := aws.NewConfig().WithRegion("us-west-1").WithCredentials(creds)
	session, err := session.NewSession(cfg)
	if err != nil {
		log.Fatal(err)
	}

	sqsService := sqs.New(session)

	randName, err := Rand256Base62()
	if err != nil {
		log.Fatal(err)
	}

	createQueueInput := &sqs.CreateQueueInput{
		QueueName: aws.String(randName), // Required
		Attributes: map[string]*string{
			sqs.QueueAttributeNameMessageRetentionPeriod: aws.String("600"),
			sqs.QueueAttributeNameVisibilityTimeout:      aws.String("0"),
		},
	}

	response, err := sqsService.CreateQueue(createQueueInput)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("created queue:", response)

	sendMessageInput := &sqs.SendMessageInput{
		MessageBody: aws.String("test"),
		QueueUrl:    response.QueueUrl,
	}

	sendResponse, err := sqsService.SendMessage(sendMessageInput)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("sent message:", sendResponse)

	receiveMessageInput := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(1),
		QueueUrl:            response.QueueUrl,
		WaitTimeSeconds:     aws.Int64(10),
	}

	receiveResponse, err := sqsService.ReceiveMessage(receiveMessageInput)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("received message:", receiveResponse)

	return
}
