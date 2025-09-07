package models

import (
	"message-scheduler/internal/domain/entity"
	"message-scheduler/internal/domain/types/status"
)

type Messages struct {
	ID              string               `gorm:"primaryKey;column:id"`
	Phone           string               `gorm:"phone"`
	Content         string               `gorm:"content"`
	Status          status.MessageStatus `gorm:"type:varchar(100);not null"`
	CreatedAt       string               `gorm:"created_at"`
	UpdatedAt       string               `gorm:"updated_at"`
	SentAt          string               `gorm:"sent_at"`
	RemoteMessageID string               `gorm:"remote_message_id"`
}

func (Messages) TableName() string {
	return "messages"
}

func MapEntityMessagesToModel(i *entity.MessagesEntity) (*Messages, error) {

	return &Messages{
		ID:              i.Id,
		Phone:           i.Phone,
		Content:         i.Content,
		Status:          i.Status,
		CreatedAt:       i.CreatedAt,
		UpdatedAt:       i.UpdatedAt,
		SentAt:          i.SentAt,
		RemoteMessageID: i.RemoteMessageId,
	}, nil
}

func MapModelMessagesToEntity(i *Messages) *entity.MessagesEntity {

	return &entity.MessagesEntity{
		Id:              i.ID,
		Phone:           i.Phone,
		Content:         i.Content,
		Status:          i.Status,
		CreatedAt:       i.CreatedAt,
		UpdatedAt:       i.UpdatedAt,
		SentAt:          i.SentAt,
		RemoteMessageId: i.RemoteMessageID,
	}
}

func MapModelMessagesToEntitySlice(messages []*Messages) []*entity.MessagesEntity {
	entities := make([]*entity.MessagesEntity, len(messages))
	for i, message := range messages {
		entities[i] = MapModelMessagesToEntity(message)
	}
	return entities
}
