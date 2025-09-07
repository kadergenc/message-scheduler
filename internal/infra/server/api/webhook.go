package api

import (
	"message-scheduler/internal/application"
	. "message-scheduler/internal/infra/server/api/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// StartSendMessageHandler godoc
// @Summary  Start Send Message
// @Description  Send a message via webhook and return the webhook response
// @Tags         messages
// @Accept       json
// @Produce      json
// @Router       /start-send-message [post]
func StartSendMessageHandler(service *application.MessageSendService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		service.StartScheduler(ctx.Context())

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
// @Router       /sent-messages [get]
func GetSentMessagesHandler(service *application.MessageSendService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
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

// StopMessageSenderHandler godoc
// @Summary  Stop Message Sender
// @Description  Stop the currently running message-sending scheduler
// @Tags         messages
// @Produce      json
// @Success      200 {object} StopSchedulerResponse "Scheduler stopped successfully"
// @Router       /stop-message-sender [post]
func StopMessageSenderHandler(service *application.MessageSendService) fiber.Handler {
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
