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
	"message-scheduler/internal/infra/database"
	"message-scheduler/internal/infra/repository"
	"message-scheduler/internal/infra/scheduler"
	"message-scheduler/internal/infra/server"
	"message-scheduler/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Logger.Info().Msg("Message Scheduler starting...")

	cfg := config.Read()

	db := database.NewPostgresDB(cfg.Postgres)

	messagesRepo := repository.NewMessagesRepository(db)

	webhookClient := webhook.NewWebhookClient(cfg.WebhookConfig.Host, time.Duration(cfg.WebhookConfig.Timeout)*time.Millisecond)

	messageScheduler := scheduler.NewSimpleScheduler()

	messageService := application.NewMessageSendService(webhookClient, messagesRepo, messageScheduler)

	appServer := server.NewAppServer(messageService)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Logger.Info().Msg("Server starting on port :" + cfg.Port)
		if err := appServer.Start(":" + cfg.Port); err != nil {
			log.Logger.Error().Err(err).Msg("Server failed to start")
		}
	}()

	log.Logger.Info().Msg("Message Scheduler started successfully. Call /start-send-message to begin processing. Press Ctrl+C to stop...")

	<-sigChan
	log.Logger.Info().Msg("Shutdown signal received. Gracefully shutting down...")

	if err := messageScheduler.Stop(); err != nil {
		log.Logger.Error().Err(err).Msg("Error stopping scheduler")
	}

	if err := appServer.Shutdown(); err != nil {
		log.Logger.Error().Err(err).Msg("Error stopping server")
	}

	log.Logger.Info().Msg("Message Scheduler stopped successfully")
}
