package application

import (
	"context"
	"fmt"
	"message-scheduler/internal/domain/entity"
	"message-scheduler/internal/domain/types/status"
	"message-scheduler/internal/infra/client/webhook"
	"message-scheduler/internal/infra/repository"
	"message-scheduler/log"
	"time"
)

type MessageSendService struct {
	client webhook.Client
	repo   repository.MessagesRepository
}

func NewMessageSendService(webhookClient webhook.Client, messagesRepo repository.MessagesRepository) *MessageSendService {
	return &MessageSendService{
		client: webhookClient,
		repo:   messagesRepo,
	}
}

func (is *MessageSendService) SendMessage(ctx context.Context, to string, content string) {

	log.Logger.Info().Str("to", to).Str("content", content).Msg("Sending message...")

	if err := is.ProcessUnsentMessages(ctx, 1); err != nil {
		log.Logger.Error().Err(err).Msg("Failed to process unsent messages")
	}
}

func (is *MessageSendService) GetUnsentMessages(ctx context.Context, limit int) ([]*entity.MessagesEntity, error) {
	log.Logger.Info().Int("limit", limit).Msg("Getting unsent messages...")

	unsentMessages, err := is.repo.GetUnsentMessages(ctx, limit)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get unsent messages")
		return nil, fmt.Errorf("failed to get unsent messages: %w", err)
	}

	log.Logger.Info().Int("count", len(unsentMessages)).Msg("Retrieved unsent messages")
	return unsentMessages, nil
}

func (is *MessageSendService) ProcessUnsentMessages(ctx context.Context, limit int) error {
	unsentMessages, err := is.GetUnsentMessages(ctx, limit)
	if err != nil {
		return err
	}

	log.Logger.Info().Int("unsent_count", len(unsentMessages)).Msg("Starting to process unsent messages")

	for _, message := range unsentMessages {
		log.Logger.Info().
			Str("message_id", message.Id).
			Str("phone", message.Phone).
			Msg("Processing unsent message")

		response, err := is.client.SendMessage(ctx, message.Phone, message.Content)
		if err != nil {
			log.Logger.Error().
				Err(err).
				Str("message_id", message.Id).
				Str("phone", message.Phone).
				Msg("Failed to send unsent message")

			message.Status = status.FAILED
			if saveErr := is.repo.Save(ctx, message); saveErr != nil {
				log.Logger.Error().Err(saveErr).Str("message_id", message.Id).Msg("Failed to update message status to FAILED")
			}
			continue
		}

		message.Status = status.SENT
		message.RemoteMessageId = response.MessageID
		message.SentAt = time.Now().Format(time.RFC3339)

		if saveErr := is.repo.Save(ctx, message); saveErr != nil {
			log.Logger.Error().Err(saveErr).Str("message_id", message.Id).Msg("Failed to update message status to SENT")
		} else {
			log.Logger.Info().
				Str("message_id", message.Id).
				Str("phone", message.Phone).
				Str("remote_message_id", message.RemoteMessageId).
				Str("sent_at", message.SentAt).
				Msg("Unsent message sent successfully and status updated")
		}
	}

	return nil
}
