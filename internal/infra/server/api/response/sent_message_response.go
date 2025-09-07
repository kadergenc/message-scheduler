package response

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
