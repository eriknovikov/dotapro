package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// getSQSClient initializes and returns an SQS client and the queue URL.
func getSQSClient(ctx context.Context) (*sqs.Client, string, error) {
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		return nil, "", fmt.Errorf("SQS_QUEUE_URL env var not set")
	}
	endpoint := os.Getenv("SQS_ENDPOINT")
	var cfg aws.Config
	var err error
	if endpoint != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: endpoint, SigningRegion: "us-east-1"}, nil
				},
			)),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to load AWS config: %w", err)
	}
	return sqs.NewFromConfig(cfg), queueURL, nil
}

// receiveMessages polls SQS for messages.
func receiveMessages(ctx context.Context, client *sqs.Client, queueURL string) (*sqs.ReceiveMessageOutput, error) {
	output, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: 10,  // Retrieve up to 10 messages at a time
		WaitTimeSeconds:     20,  // Enable long polling
		VisibilityTimeout:   300, // 5 minutes visibility timeout
	})
	if err != nil {
		return nil, fmt.Errorf("error receiving SQS messages: %w", err)
	}
	return output, nil
}

// deleteMessage deletes a message from the SQS queue.
func deleteMessage(ctx context.Context, client *sqs.Client, queueURL string, receiptHandle *string) error {
	_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		return fmt.Errorf("error deleting message from SQS: %w", err)
	}
	return nil
}
