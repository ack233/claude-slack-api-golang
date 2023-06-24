package slack

import (
	"fmt"
	"slackapi/pkgs/zlog"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func middlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	zlog.SugLog.Info("Connecting to Slack with Socket Mode...")
}

func middlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	zlog.SugLog.Info("Connection failed. Retrying later...")
}

func middlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	zlog.SugLog.Info("Connected to Slack with Socket Mode.")
}

type MessageEvent struct {
	// Basic Message Event - https://api.slack.com/events/message
	ClientMsgID     string `json:"client_msg_id"`
	Type            string `json:"type"`
	User            string `json:"user"`
	Text            string `json:"text"`
	ThreadTimeStamp string `json:"thread_ts"`
	TimeStamp       string `json:"ts"`
	Channel         string `json:"channel"`
	ChannelType     string `json:"channel_type"`
	EventTimeStamp  string `json:"event_ts"`

	// When Message comes from a channel that is shared between workspaces
	UserTeam   string `json:"user_team,omitempty"`
	SourceTeam string `json:"source_team,omitempty"`

	// Edited Message
	Message         *MessageEvent `json:"message,omitempty"`
	PreviousMessage *MessageEvent `json:"previous_message,omitempty"`
	Edited          *Edited       `json:"edited,omitempty"`

	// Message Subtypes
	SubType string `json:"subtype,omitempty"`

	// bot_message (https://api.slack.com/events/message/bot_message)
	BotID    string `json:"bot_id,omitempty"`
	Username string `json:"username,omitempty"`
	//Icons    *Icon  `json:"icons,omitempty"`

	Upload bool `json:"upload"`
	//Files  []File `json:"files"`

	Attachments []slack.Attachment `json:"attachments,omitempty"`

	// Root is the message that was broadcast to the channel when the SubType is
	// thread_broadcast. If this is not a thread_broadcast message event, this
	// value is nil.
	Root *MessageEvent `json:"root"`
}

func middlewareMessageEvent(evt *socketmode.Event, client *socketmode.Client) {
	client.Ack(*evt.Request)

	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)

	if !ok {
		zlog.SugLog.Errorf("Unexpected event data type")
		return
	}

	msgEvent, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if !ok {
		zlog.SugLog.Errorf("Unexpected inner event data type")
		return
	}

	if msgEvent.Message != nil && msgEvent.SubType == "message_changed" {

		threadTS := msgEvent.Message.ThreadTimeStamp
		text := msgEvent.Message.Text
		zlog.SugLog.Debugf("get event: %v", threadTS)

		ResponseChannels.SendMessage(threadTS, text)

		// Process the message event here...

	}

}

func middlewareSlashCommand(evt *socketmode.Event, client *socketmode.Client) {
	cmd, ok := evt.Data.(slack.SlashCommand)
	if !ok {
		zlog.SugLog.Warnf("Ignored %+v\n", evt)
		return
	}

	client.Debugf("Slash command received: %+v", cmd)

	client.Ack(*evt.Request)
	responseText := fmt.Sprintf("Hi, <@%s>!", cmd.UserID)

	_, _, err := client.Client.PostMessage(cmd.ChannelID, slack.MsgOptionText(responseText, false))

	if err != nil {
		zlog.SugLog.Errorf("Failed posting message: %v", err)
	}
}
