package entity

import "message-scheduler/internal/domain/types/status"

type MessagesEntity struct {
	Id              string
	Phone           string
	Content         string
	Status          status.MessageStatus
	CreatedAt       string
	UpdatedAt       string
	SentAt          string
	RemoteMessageId string
}
