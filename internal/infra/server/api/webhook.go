package api

import (
	"message-scheduler/internal/application"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MessageRequest struct {
	To      string `json:"to" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type MessageResponse struct {
	Message   string `json:"message" example:"accepted"`
	MessageID string `json:"messageId" example:"01623bff-7fa9-4ccb-a843-e6d98908dc49"`
}

type SentMessageResponse struct {
	ID              string `json:"id" example:"01623bff-7fa9-4ccb-a843-e6d98908dc49"`
	Phone           string `json:"phone" example:"+905551234567"`
	Content         string `json:"content" example:"Hello, World!"`
	Status          string `json:"status" example:"SENT"`
	CreatedAt       string `json:"createdAt" example:"2023-10-01T10:00:00Z"`
	UpdatedAt       string `json:"updatedAt" example:"2023-10-01T10:05:00Z"`
	SentAt          string `json:"sentAt" example:"2023-10-01T10:05:00Z"`
	RemoteMessageID string `json:"remoteMessageId" example:"whatsapp-msg-123"`
}

type GetSentMessagesResponse struct {
	Messages []SentMessageResponse `json:"messages"`
	Total    int                   `json:"total"`
}

type StopSchedulerResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Message scheduler stopped successfully"`
}

// MessageHandler godoc
// @Summary  Start Send Message
// @Description  Send a message via webhook and return the webhook response
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        request body MessageRequest true "Message request payload"
// @Success      200 {object} MessageResponse "Message sent successfully"
// @Failure      400 {object} map[string]string "Invalid request body or missing required fields"
// @Failure      500 {object} map[string]string "Failed to send message"
// @Router       /start-send-message [post]
func MessageHandler(service *application.MessageSendService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body MessageRequest

		if err := ctx.BodyParser(&body); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if body.To == "" || body.Content == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Both 'to' and 'content' fields are required"})
		}

		service.SendMessage(ctx.Context(), body.To, body.Content)

		return nil
	}
}

// GetSentMessagesHandler godoc
// @Summary  Get Sent Messages
// @Description  Retrieve sent messages with optional limit
// @Tags         messages
// @Produce      json
// @Param        limit query int false "Number of messages to retrieve (default: 50, max: 100)"
// @Success      200 {object} GetSentMessagesResponse "Sent messages retrieved successfully"
// @Failure      400 {object} map[string]string "Invalid limit parameter"
// @Failure      500 {object} map[string]string "Failed to retrieve sent messages"
// @Router       /sent-messages [get]
func GetSentMessagesHandler(service *application.MessageSendService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Default limit is 50, maximum is 100
		limit := 50
		if limitParam := ctx.Query("limit"); limitParam != "" {
			parsedLimit, err := strconv.Atoi(limitParam)
			if err != nil || parsedLimit <= 0 {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit parameter. Must be a positive integer"})
			}
			if parsedLimit > 100 {
				parsedLimit = 100
			}
			limit = parsedLimit
		}

		sentMessages, err := service.GetSentMessages(ctx.Context(), limit)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve sent messages"})
		}

		responseMessages := make([]SentMessageResponse, len(sentMessages))
		for i, message := range sentMessages {
			responseMessages[i] = SentMessageResponse{
				ID:              message.Id,
				Phone:           message.Phone,
				Content:         message.Content,
				Status:          string(message.Status),
				CreatedAt:       message.CreatedAt,
				UpdatedAt:       message.UpdatedAt,
				SentAt:          message.SentAt,
				RemoteMessageID: message.RemoteMessageId,
			}
		}

		response := GetSentMessagesResponse{
			Messages: responseMessages,
			Total:    len(responseMessages),
		}

		return ctx.JSON(response)
	}
}

// StopSchedulerHandler godoc
// @Summary  Stop Message Sender
// @Description  Stop the currently running message-sending scheduler
// @Tags         messages
// @Produce      json
// @Success      200 {object} StopSchedulerResponse "Scheduler stopped successfully"
// @Failure      500 {object} map[string]string "Failed to stop scheduler"
// @Router       /stop-message-sender [post]
func StopSchedulerHandler(service *application.MessageSendService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := service.StopScheduler()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to stop scheduler: " + err.Error()})
		}

		response := StopSchedulerResponse{
			Status:  "success",
			Message: "Message scheduler stopped successfully",
		}

		return ctx.JSON(response)
	}
}
