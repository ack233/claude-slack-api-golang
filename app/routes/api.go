package routes

import (

	//"filrserver/pkgs/zlog"

	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"slackapi/app/middlewares"
	"slackapi/pkgs/requests"
	"slackapi/pkgs/zlog"
	"slackapi/slack"

	"slackapi/pkgs/config"

	"github.com/gin-gonic/gin"
	"github.com/kyokomi/emoji/v2"
)

func login(c *gin.Context) {
	clientID := config.SlackConfig.SLACK_OAUTH_CLIENT_ID
	redirectURI := config.SlackConfig.SLACK_OAUTH_REDIRECT_URI

	url := fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?"+
			"user_scope=chat:write&"+
			"scope=chat:write,users:read,channels:history&"+
			"client_id=%s&"+
			"redirect_uri=%s",
		clientID,
		redirectURI,
	)

	c.Redirect(http.StatusMovedPermanently, url)
}

func callback(c *gin.Context) {
	error := c.DefaultQuery("error", "")
	if error != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": fmt.Sprintf("OAuth error: %s", error),
		})
		return
	}

	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "Missing OAuth code",
		})
		return
	}

	clientID := config.SlackConfig.SLACK_OAUTH_CLIENT_ID
	clientSecret := config.SlackConfig.SLACK_OAUTH_CLIENT_SECRET
	redirectURI := config.SlackConfig.SLACK_OAUTH_REDIRECT_URI

	var result SlackOAuthResponse

	_, err := requests.Client.R().
		//SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(
			map[string]string{
				"client_id":     clientID,
				"client_secret": clientSecret,
				"code":          code,
				"redirect_uri":  redirectURI,
			}).
		SetResult(&result).
		Post("https://slack.com/api/oauth.v2.access")

	if result.Ok {

		token, err := middlewares.CreateToken(result.AuthedUser.AccessToken)
		zlog.Errorerror(err)
		c.HTML(http.StatusOK, "success.html", gin.H{
			"access_token": token,
		})
	} else {
		zlog.Errorerror(err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error_msg": result.Error,
		})
	}

}

func revoke(c *gin.Context) {

	accessToken := c.GetString("token")
	resp, err := requests.Client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Get("https://slack.com/api/auth.revoke")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, resp.Header().Get("Content-Type"), resp.Body())
}

func conversation(c *gin.Context) {
	var req ConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken := c.GetString("token")

	channel := req.ChannelID

	prompt := strings.Join(req.Messages[0].Content.Parts, "")
	botID := "bot"
	channelParts := strings.Split(channel, ":")
	if len(channelParts) > 1 {
		channel = channelParts[0]
		botID = channelParts[1]
	}

	payload := map[string]string{
		"text":       fmt.Sprintf("<@%s> %s", botID, prompt),
		"channel":    channel,
		"thread_ts":  req.ConversationID,
		"link_names": "true",
	}

	slackResponse := SlackResponse{}
	_, err := requests.Client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetFormData(payload).
		SetResult(&slackResponse).
		Post("https://slack.com/api/chat.postMessage")

	if err != nil {
		c.JSON(http.StatusBadRequest, ConversationResponse{Error: err.Error()})
		zlog.SugLog.Error(err)
		return
	}

	if slackResponse.Error != "" {
		zlog.SugLog.Error(slackResponse.Error)

		c.JSON(http.StatusBadRequest, ConversationResponse{Error: slackResponse.Error})
		return
	}

	userTs := req.ConversationID
	if userTs == "" {
		userTs = slackResponse.Message.Ts
	}
	zlog.SugLog.Debugf("start with userTs: %s", userTs)
	ch := slack.ResponseChannels.CreateChannel(userTs)

	c.SSEvent("ping", "") // 在c.Stream之前发送ping事件
	zlog.SugLog.Debug("test")
	c.Stream(func(w io.Writer) bool {
		return emitSSE(c, ch, userTs)
	})

}

func emitSSE(c *gin.Context, ch chan string, key string) (enevtFlag bool) {

	defer func() {
		if !enevtFlag {
			slack.ResponseChannels.DeleteChannel(key)
		}
	}()

	var enevtDone bool
	var typingSuffix string = "_Typing…_"

	select {
	case msg, ok := <-ch:
		if !ok {
			zlog.SugLog.Debug("event channel not ok")
			return
		}

		msg = parseMsg(msg)

		zlog.SugLog.Debugf("get message: %s", msg)

		if strings.HasSuffix(msg, typingSuffix) {
			msg = strings.TrimSuffix(msg, typingSuffix)
			enevtFlag = true

		} else {
			enevtDone = true
		}

		response := ConversationResponse{
			Message: Message{
				ID:      "someId",
				Role:    "assistant",
				Content: Content{ContentType: "text", Parts: []string{msg}},
				Author: struct {
					Role string `json:"role"`
				}{Role: "assistant"},
			},
			ConversationID: key,
		}

		data, err := json.Marshal(response)
		zlog.Errorerror(err)
		c.SSEvent("data", string(data))

	case <-time.After(7 * time.Second):
		c.SSEvent("data", "[TimeOutError]")
		zlog.SugLog.Debug("event channel TimeOutError")
	}

	if enevtDone {
		c.SSEvent("data", "[DONE]")
		zlog.SugLog.Debug("event channel done")
	}
	return

}

func parseMsg(msg string) string {

	msg = strings.TrimSpace(msg)
	msg = emoji.Sprint(msg)
	msg = html.UnescapeString(msg)
	//msg = strconv.QuoteToASCII(msg)
	return msg

}
