package telegram

type (
	UpdateResponse struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	Update struct {
		Id      int              `json:"update_id"`
		Message *IncomingMessage `json:"message"`
	}

	IncomingMessage struct {
		Text string `json:"text"`
		From From   `json:"from"`
		Chat Chat   `json:"chat"`
	}

	From struct {
		Username string `json:"username"`
	}

	Chat struct {
		Id int `json:"id"`
	}
)
