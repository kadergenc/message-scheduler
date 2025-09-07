package repository

import (
	"context"
	"fmt"
	"message-scheduler/internal/domain/entity"
	"message-scheduler/internal/infra/repository/models"
	"message-scheduler/log"

	"gorm.io/gorm"
)

type MessagesRepository interface {
	Save(ctx context.Context, message *entity.MessagesEntity) error
	GetUnsentMessages(ctx context.Context, recordLimit int) ([]*entity.MessagesEntity, error)
}

type PostgresMessagesRepository struct {
	db *gorm.DB
}

func NewMessagesRepository(db *gorm.DB) *PostgresMessagesRepository {
	db.Logger = &GormLogger{log.Logger}

	return &PostgresMessagesRepository{db: db}
}

func (r *PostgresMessagesRepository) Save(ctx context.Context, i *entity.MessagesEntity) error {
	message, err := models.MapEntityMessagesToModel(i)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(&message).Error; err != nil {
		log.Logger.Error().Msgf("Error saving message: %+v", err)
		return fmt.Errorf("failed to save message with id=%s: %w", message.ID, err)
	}

	log.Logger.Info().Str("messageId", message.ID).Str("status", string(message.Status)).Msg("saved message")
	return nil
}

func (r *PostgresMessagesRepository) GetUnsentMessages(ctx context.Context, recordLimit int) ([]*entity.MessagesEntity, error) {
	var messages []*models.Messages

	err := r.db.WithContext(ctx).
		Where("status = ?", "unsent").
		Limit(recordLimit).
		Find(&messages).Error

	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to fetch unsent messages from database")
		return nil, fmt.Errorf("failed to fetch unsent messages: %w", err)
	}

	log.Logger.Info().Int("found_messages", len(messages)).Msg("Retrieved unsent messages from database")

	log.Logger.Debug().Interface("messages", messages).Msg("Retrieved messages content")

	return models.MapModelMessagesToEntitySlice(messages), nil
}
