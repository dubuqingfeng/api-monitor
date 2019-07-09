package senders

import (
	"bytes"
	"fmt"
	"github.com/dubuqingfeng/api-monitor/models"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type SlackSender struct {
	Sender
}

type SlackMessage struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	AsUser  bool   `json:"as_user"`
	Text    string `json:"text"`
	Token   string `json:"token"`
}

var SlackPusher SlackSender

func init() {
	SlackPusher = SlackSender{}
}

func (s SlackSender) IsSupport() bool {
	return utils.Config.SenderConfig.Slack.IsEnabled
}

func (s SlackSender) Send(notifications []*models.Notification) {
	if !utils.Config.SenderConfig.Slack.IsEnabled {
		return
	}
	for _, item := range notifications {
		s.SingleSend(item)
	}
}

func (s SlackSender) SingleSend(notification *models.Notification) {
	message := SlackMessage{
		AsUser:  true,
		Channel: utils.Config.SenderConfig.Slack.Channel,
		Text:    s.BuildMessage(notification),
	}
	data := url.Values{}
	data.Set("token", utils.Config.SenderConfig.Slack.RobotToken)
	data.Add("channel", message.Channel)
	data.Add("text", message.Text)
	data.Add("as_user", strconv.FormatBool(message.AsUser))

	body, err := http.Post("https://slack.com/api/chat.postMessage", "application/x-www-form-urlencoded",
		bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Error(err)
	}
	content, err := ioutil.ReadAll(body.Body)
	log.Info(string(content))
}

func (s SlackSender) BuildMessage(notification *models.Notification) string {
	return fmt.Sprintf("status:%d,type:%s,url:%s", notification.HttpStatus, notification.Reason, notification.URL)
}
