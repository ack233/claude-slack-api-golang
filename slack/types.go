package slack

type MessageChangeEvent struct {
	Token               string          `json:"token"`
	TeamID              string          `json:"team_id"`
	ContextTeamID       string          `json:"context_team_id"`
	ContextEnterpriseID interface{}     `json:"context_enterprise_id"`
	APIAppID            string          `json:"api_app_id"`
	Event               Event           `json:"event"`
	Type                string          `json:"type"`
	EventID             string          `json:"event_id"`
	EventTime           int64           `json:"event_time"`
	Authorizations      []Authorization `json:"authorizations"`
	IsExtSharedChannel  bool            `json:"is_ext_shared_channel"`
	EventContext        string          `json:"event_context"`
}

type Event struct {
	Type    string  `json:"type"`
	Subtype string  `json:"subtype"`
	Message Message `json:"message"`
}

type Message struct {
	BotID        string     `json:"bot_id"`
	Type         string     `json:"type"`
	Text         string     `json:"text"`
	User         string     `json:"user"`
	AppID        string     `json:"app_id"`
	Blocks       []Block    `json:"blocks"`
	Team         string     `json:"team"`
	BotProfile   BotProfile `json:"bot_profile"`
	Edited       Edited     `json:"edited"`
	ThreadTS     string     `json:"thread_ts"`
	ParentUserID string     `json:"parent_user_id"`
	Ts           string     `json:"ts"`
	SourceTeam   string     `json:"source_team"`
	UserTeam     string     `json:"user_team"`
}

type Block struct {
	Type     string    `json:"type"`
	BlockID  string    `json:"block_id"`
	Elements []Element `json:"elements"`
}

type Element struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Style Style  `json:"style,omitempty"`
}

type Style struct {
	Italic bool `json:"italic"`
}

type BotProfile struct {
	ID      string `json:"id"`
	AppID   string `json:"app_id"`
	Name    string `json:"name"`
	Icons   Icons  `json:"icons"`
	Deleted bool   `json:"deleted"`
	Updated int64  `json:"updated"`
	TeamID  string `json:"team_id"`
}

type Icons struct {
	Image36 string `json:"image_36"`
	Image48 string `json:"image_48"`
	Image72 string `json:"image_72"`
}

type Edited struct {
	User string `json:"user"`
	Ts   string `json:"ts"`
}

type Authorization struct {
	EnterpriseID        interface{} `json:"enterprise_id"`
	TeamID              string      `json:"team_id"`
	UserID              string      `json:"user_id"`
	IsBot               bool        `json:"is_bot"`
	IsEnterpriseInstall bool        `json:"is_enterprise_install"`
}
