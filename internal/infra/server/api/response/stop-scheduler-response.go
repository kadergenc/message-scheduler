package response

type StopSchedulerResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Message scheduler stopped successfully"`
}
