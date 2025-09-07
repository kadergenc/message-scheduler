// @title Message Scheduler API
// @version 1.0
// @description This is a message scheduler service
// @host localhost:8081
// @BasePath /
//
//go:generate swag init
package main

import (
	"message-scheduler/config"
	_ "message-scheduler/docs"
	"message-scheduler/internal/application"
	"message-scheduler/internal/infra/client/webhook"
	"message-scheduler/internal/infra/server"
	"message-scheduler/log"
	"time"
)

func main() {
	log.Logger.Info().Msg("Message Scheduler starting...")

	cfg := config.Read()

	webhookClient := webhook.NewWebhookClient("http://localhost:8000/webhook", 30*time.Second)

	messageService := application.NewMessageSendService(*webhookClient)

	appServer := server.NewAppServer(messageService)

	log.Logger.Info().Msg("Server starting on port :" + cfg.Port)
	if err := appServer.Start(":" + cfg.Port); err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
