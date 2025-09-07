package status

type MessageStatus string

const (
	UNSENT MessageStatus = "UNSENT"
	SENT   MessageStatus = "SENT"
	FAILED MessageStatus = "FAILED"
)
