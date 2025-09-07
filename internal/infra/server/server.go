package server

import (
	"message-scheduler/internal/application"
	"message-scheduler/internal/infra/server/api"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Message scheduler alive!"`
}

type AppServer struct {
	app     *fiber.App
	service *application.MessageSendService
}

func NewAppServer(service *application.MessageSendService) AppServer {
	app := fiber.New()

	app.Post("/start-send-message", api.StartSendMessageHandler(service))
	app.Post("/stop-message-sender", api.StopMessageSenderHandler(service))
	app.Get("/sent-messages", api.GetSentMessagesHandler(service))
	app.Get("/_monitoring/health", index)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Use(func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/")
	})

	return AppServer{app: app, service: service}
}

func (a AppServer) Start(port string) error {
	return a.app.Listen(port)
}

func (a AppServer) Shutdown() error {
	return a.app.Shutdown()
}

// HealthCheck godoc
// @Summary  Health Check
// @Description  Check if the message scheduler service is running
// @Tags         monitoring
// @Produce      json
// @Success      200 {object} HealthResponse "Service is healthy"
// @Router       /_monitoring/health [get]
func index(ctx *fiber.Ctx) error {
	return ctx.JSON(HealthResponse{
		Status:  "ok",
		Message: "Message scheduler alive!",
	})
}
