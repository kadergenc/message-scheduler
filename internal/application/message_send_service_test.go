package application

import (
	"context"
	"fmt"
	"message-scheduler/internal/domain/entity"
	"message-scheduler/internal/domain/types/status"
	"message-scheduler/internal/infra/client/webhook"
	"message-scheduler/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestMessage(messageStatus status.MessageStatus) *entity.MessagesEntity {
	return &entity.MessagesEntity{
		Id:              uuid.New().String(),
		Phone:           "+905551234567",
		Content:         "Test message content",
		Status:          messageStatus,
		CreatedAt:       time.Now().Format(time.RFC3339),
		UpdatedAt:       time.Now().Format(time.RFC3339),
		SentAt:          "",
		RemoteMessageId: "",
	}
}

func TestSendMessage_Success(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	mockScheduler.On("ScheduleJob", mock.Anything, 2*time.Minute).Return()
	mockScheduler.On("Start", mock.Anything).Return()

	ctx := context.Background()
	service.SendMessage(ctx, "+905551234567", "Test message")

	time.Sleep(10 * time.Millisecond)

	assert.True(t, service.schedulerRunning)
	mockScheduler.AssertExpectations(t)
}

func TestProcessUnsentMessages_Success(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	ctx := context.Background()
	limit := 5

	unsentMessages := []*entity.MessagesEntity{
		createTestMessage(status.UNSENT),
		createTestMessage(status.UNSENT),
	}

	webhookResponse := &webhook.WebhookResponse{
		Message:   "Message sent successfully",
		MessageID: "webhook-msg-123",
	}

	mockRepo.On("GetUnsentMessages", ctx, limit).Return(unsentMessages, nil)
	mockWebhook.On("SendMessage", ctx, unsentMessages[0].Phone, unsentMessages[0].Content).Return(webhookResponse, nil)
	mockWebhook.On("SendMessage", ctx, unsentMessages[1].Phone, unsentMessages[1].Content).Return(webhookResponse, nil)
	mockRepo.On("Save", ctx, mock.MatchedBy(func(msg *entity.MessagesEntity) bool {
		return msg.Status == status.SENT && msg.RemoteMessageId == "webhook-msg-123"
	})).Return(nil).Twice()

	err := service.ProcessUnsentMessages(ctx, limit)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
}

func TestProcessUnsentMessages_WebhookFailure(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	ctx := context.Background()
	limit := 2

	unsentMessages := []*entity.MessagesEntity{
		createTestMessage(status.UNSENT),
	}

	mockRepo.On("GetUnsentMessages", ctx, limit).Return(unsentMessages, nil)
	mockWebhook.On("SendMessage", ctx, unsentMessages[0].Phone, unsentMessages[0].Content).Return(nil, fmt.Errorf("webhook error"))
	mockRepo.On("Save", ctx, mock.MatchedBy(func(msg *entity.MessagesEntity) bool {
		return msg.Status == status.FAILED
	})).Return(nil)

	err := service.ProcessUnsentMessages(ctx, limit)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
}

func TestProcessUnsentMessages_SaveFailure(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	ctx := context.Background()
	limit := 1

	unsentMessages := []*entity.MessagesEntity{
		createTestMessage(status.UNSENT),
	}

	webhookResponse := &webhook.WebhookResponse{
		Message:   "Message sent successfully",
		MessageID: "webhook-msg-123",
	}

	mockRepo.On("GetUnsentMessages", ctx, limit).Return(unsentMessages, nil)
	mockWebhook.On("SendMessage", ctx, unsentMessages[0].Phone, unsentMessages[0].Content).Return(webhookResponse, nil)
	mockRepo.On("Save", ctx, mock.Anything).Return(fmt.Errorf("database error"))

	err := service.ProcessUnsentMessages(ctx, limit)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
}

func TestGetUnsentMessages_Success(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	ctx := context.Background()
	limit := 10

	expectedMessages := []*entity.MessagesEntity{
		createTestMessage(status.UNSENT),
		createTestMessage(status.UNSENT),
	}

	mockRepo.On("GetUnsentMessages", ctx, limit).Return(expectedMessages, nil)

	messages, err := service.GetUnsentMessages(ctx, limit)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	assert.Len(t, messages, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetUnsentMessages_RepositoryError(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	ctx := context.Background()
	limit := 10

	mockRepo.On("GetUnsentMessages", ctx, limit).Return(nil, fmt.Errorf("database connection failed"))

	messages, err := service.GetUnsentMessages(ctx, limit)

	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.Contains(t, err.Error(), "failed to get unsent messages")
	mockRepo.AssertExpectations(t)
}

func TestGetSentMessages_Success(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	ctx := context.Background()
	limit := 5

	sentMessage := createTestMessage(status.SENT)
	sentMessage.SentAt = time.Now().Format(time.RFC3339)
	sentMessage.RemoteMessageId = "webhook-123"

	expectedMessages := []*entity.MessagesEntity{sentMessage}

	mockRepo.On("GetSentMessages", ctx, limit).Return(expectedMessages, nil)

	messages, err := service.GetSentMessages(ctx, limit)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	assert.Len(t, messages, 1)
	assert.Equal(t, status.SENT, messages[0].Status)
	mockRepo.AssertExpectations(t)
}

func TestGetSentMessages_RepositoryError(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	ctx := context.Background()
	limit := 5

	mockRepo.On("GetSentMessages", ctx, limit).Return(nil, fmt.Errorf("database timeout"))

	messages, err := service.GetSentMessages(ctx, limit)

	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.Contains(t, err.Error(), "failed to get sent messages")
	mockRepo.AssertExpectations(t)
}

func TestStopScheduler_Success(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	service.schedulerRunning = true

	mockScheduler.On("Stop").Return(nil)

	err := service.StopScheduler()

	assert.NoError(t, err)
	assert.False(t, service.schedulerRunning)
	mockScheduler.AssertExpectations(t)
}

func TestStopScheduler_NotRunning(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	assert.False(t, service.schedulerRunning)

	err := service.StopScheduler()

	assert.NoError(t, err)
	assert.False(t, service.schedulerRunning)

}

func TestStopScheduler_NotInitialized(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, nil)

	err := service.StopScheduler()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scheduler is not initialized")
}

func TestStopScheduler_StopError(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	service.schedulerRunning = true

	mockScheduler.On("Stop").Return(fmt.Errorf("failed to stop scheduler"))

	err := service.StopScheduler()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stop scheduler")
	assert.True(t, service.schedulerRunning) // Should remain true on error
	mockScheduler.AssertExpectations(t)
}

func TestContinuousMessageProcessorJob_Execute(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	job := &continuousMessageProcessorJob{
		messageService: service,
		limit:          3,
	}

	ctx := context.Background()

	mockRepo.On("GetUnsentMessages", ctx, 3).Return([]*entity.MessagesEntity{}, nil)

	err := job.Execute(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "ContinuousMessageProcessor", job.Name())
	mockRepo.AssertExpectations(t)
}

func TestSendMessage_SchedulerAlreadyRunning(t *testing.T) {
	mockWebhook := &mocks.WebhookClientMock{}
	mockRepo := &mocks.MessagesRepositoryMock{}
	mockScheduler := &mocks.SchedulerMock{}

	service := NewMessageSendService(mockWebhook, mockRepo, mockScheduler)

	service.schedulerRunning = true

	ctx := context.Background()
	service.SendMessage(ctx, "+905551234567", "Test message")

	assert.True(t, service.schedulerRunning)
}
