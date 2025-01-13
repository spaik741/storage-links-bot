package telegram

type (
	UpdateResponse struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	Update struct {
		Id      int    `json:"update_id"`
		Message string `json:"message"`
	}
)
