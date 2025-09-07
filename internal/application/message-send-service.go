package application

import (
	"context"
	"fmt"
	"message-scheduler/internal/infra/client/webhook"
	"message-scheduler/log"
)

type MessageSendService struct {
	client webhook.Client
}

func NewMessageSendService(webhookClient webhook.Client) *MessageSendService {
	return &MessageSendService{
		client: webhookClient,
	}
}

func (is *MessageSendService) SendMessage(ctx context.Context, to string, content string) (*webhook.WebhookResponse, error) {

	log.Logger.Info().Str("to", to).Str("content", content).Msg("Sending message...")

	response, err := is.client.SendMessage(ctx, to, content)
	if err != nil {
		log.Logger.Error().Err(err).Str("to", to).Msg("Failed to send message")
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	log.Logger.Info().Str("to", to).Str("messageId", response.MessageID).Msg("Message sent successfully")
	return response, nil
}
