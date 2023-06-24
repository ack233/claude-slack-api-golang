package slack

import (
	"log"
	"os"

	"github.com/slack-go/slack/slackevents"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	"slackapi/pkgs/config"
	"slackapi/pkgs/zlog"
)

type Client struct {
	SlackClient       *slack.Client
	SocketModeClient  *socketmode.Client
	SocketModeHandler *socketmode.SocketmodeHandler
}

var ResponseChannels *ChannelManager = new(ChannelManager)

func NewClient() *Client {
	appToken := config.SlackConfig.SLACK_APP_TOKEN
	botToken := config.SlackConfig.SLACK_BOT_TOKEN

	slackClient := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	socketModeClient := socketmode.New(slackClient,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	socketModeHandler := socketmode.NewSocketmodeHandler(socketModeClient)

	client := &Client{
		SlackClient:       slackClient,
		SocketModeClient:  socketModeClient,
		SocketModeHandler: socketModeHandler,
	}

	return client
}

func (c *Client) InitSlackEventHandlers() {
	c.SocketModeHandler.Handle(socketmode.EventTypeConnecting, middlewareConnecting)
	c.SocketModeHandler.Handle(socketmode.EventTypeConnectionError, middlewareConnectionError)
	c.SocketModeHandler.Handle(socketmode.EventTypeConnected, middlewareConnected)

	//c.SocketModeHandler.Handle(socketmode.EventTypeEventsAPI, middlewareEventsAPI)

	c.SocketModeHandler.HandleEvents(slackevents.Message, middlewareMessageEvent)
	c.SocketModeHandler.HandleSlashCommand("/hello-socket-mode", middlewareSlashCommand)

}

func Run() {
	c := NewClient()
	c.InitSlackEventHandlers()
	err := c.SocketModeHandler.RunEventLoop()
	//err := c.SocketModeClient.Run()
	zlog.Errorerror(err)
}
