package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Logger()
}

func main() {

	ctx := context.Background()

	pool, err := getPostgresClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer pool.Close()
	log.Info().Msg("PostgreSQL client initialized.")

	sqsClient, queueURL, err := getSQSClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create SQS client")
	}
	log.Info().Msg("SQS client initialized.")
	for {
		log.Info().Msg("Polling SQS for messages...")
		output, err := receiveMessages(ctx, sqsClient, queueURL)
		if err != nil {
			log.Error().Err(err).Msg("Error receiving SQS messages")
			time.Sleep(10 * time.Second) // Wait before retrying
			continue
		}

		if len(output.Messages) == 0 {
			log.Info().Msg("No messages received. Waiting...")
			time.Sleep(10 * time.Second) // Wait before next poll
			continue
		}

		for _, message := range output.Messages {
			log.Info().Str("message_id", *message.MessageId).Msg("Received message")
			if err := processMessage(ctx, pool, *message.Body); err != nil {
				log.Error().Err(err).Str("message_id", *message.MessageId).Msg("Error processing message")
				// Message will become visible again after VisibilityTimeout for retry
				continue
			}

			// Delete the message from the queue after successful processing

			if err := deleteMessage(ctx, sqsClient, queueURL, message.ReceiptHandle); err != nil {
				log.Error().Err(err).Str("message_id", *message.MessageId).Msg("Error deleting message from SQS")
			} else {
				log.Info().Str("message_id", *message.MessageId).Msg("Successfully processed and deleted message")
			}
		}
	}
}
