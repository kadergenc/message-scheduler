package application

import (
	"context"
	"fmt"
	"message-scheduler/internal/domain/entity"
	"message-scheduler/internal/domain/types/status"
	"message-scheduler/internal/infra/client/webhook"
	"message-scheduler/internal/infra/repository"
	"message-scheduler/internal/port"
	"message-scheduler/log"
	"time"
)

type MessageSendService struct {
	client           webhook.WebhookClient
	repo             repository.MessagesRepository
	scheduler        port.Scheduler
	schedulerRunning bool
}

func NewMessageSendService(webhookClient webhook.WebhookClient, messagesRepo repository.MessagesRepository, scheduler port.Scheduler) *MessageSendService {
	return &MessageSendService{
		client:           webhookClient,
		repo:             messagesRepo,
		scheduler:        scheduler,
		schedulerRunning: false,
	}
}

func (is *MessageSendService) StartScheduler(ctx context.Context) {

	if !is.schedulerRunning && is.scheduler != nil {
		log.Logger.Info().Msg("Starting continuous message processing scheduler...")

		continuousJob := &continuousMessageProcessorJob{
			messageService: is,
			limit:          2,
		}

		is.scheduler.ScheduleJob(continuousJob, 2*time.Minute)

		go func() {
			is.scheduler.Start(ctx)
		}()

		is.schedulerRunning = true
		log.Logger.Info().Msg("Continuous message processing scheduler started successfully")
	}

}

func (is *MessageSendService) StopScheduler() error {
	if is.scheduler == nil {
		return fmt.Errorf("scheduler is not initialized")
	}

	if !is.schedulerRunning {
		log.Logger.Warn().Msg("Scheduler is not running, ignoring stop request")
		return nil
	}

	log.Logger.Info().Msg("Stopping message scheduler via API request")
	err := is.scheduler.Stop()
	if err == nil {
		is.schedulerRunning = false
		log.Logger.Info().Msg("Message scheduler stopped successfully")
	}
	return err
}

type continuousMessageProcessorJob struct {
	messageService *MessageSendService
	limit          int
}

func (j *continuousMessageProcessorJob) Execute(ctx context.Context) error {
	return j.messageService.ProcessUnsentMessages(ctx, j.limit)
}

func (j *continuousMessageProcessorJob) Name() string {
	return "ContinuousMessageProcessor"
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

func (is *MessageSendService) GetSentMessages(ctx context.Context, limit int) ([]*entity.MessagesEntity, error) {
	log.Logger.Info().Int("limit", limit).Msg("Getting sent messages...")

	sentMessages, err := is.repo.GetSentMessages(ctx, limit)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get sent messages")
		return nil, fmt.Errorf("failed to get sent messages: %w", err)
	}

	log.Logger.Info().Int("count", len(sentMessages)).Msg("Retrieved sent messages")
	return sentMessages, nil
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
