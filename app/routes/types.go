package routes

type SlackOAuthResponse struct {
	Ok         bool            `json:"ok"`
	AuthedUser SlackAuthedUser `json:"authed_user"`
	Error      string          `json:"error"`
}

type SlackAuthedUser struct {
	AccessToken string `json:"access_token"`
}

type ConversationRequest struct {
	Messages       []Message `json:"messages"`
	ConversationID string    `json:"conversation_id"`
	ChannelID      string    `json:"channel_id"`
}

type Message struct {
	Role    string  `json:"role"`
	Content Content `json:"content"`
	Author  Author  `json:"author"`
	ID      string  `json:"id"`
}

type Content struct {
	ContentType string `json:"content_type"`
	Parts       []string `json:"parts"`
}

type Author struct {
	Role string `json:"role"`
}

type ConversationResponse struct {
	Message        Message `json:"message"`
	ConversationID string  `json:"conversation_id"`
	Error          string  `json:"error"`
}

type SlackResponse struct {
	Message struct {
		User string `json:"user"`
		Ts   string `json:"ts"`
	} `json:"message"`
	Error string `json:"error"`
	Ok    bool   `json:"ok"`
}
