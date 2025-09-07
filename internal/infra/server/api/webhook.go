package api

import (
	"message-scheduler/internal/application"

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
