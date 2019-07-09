package senders

import (
	"fmt"
	"github.com/bearyinnovative/bearychat-go"
	"github.com/dubuqingfeng/api-monitor/models"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type BearyChatSender struct {
	Sender
}

var BearyChatPusher BearyChatSender

func init() {
	BearyChatPusher = BearyChatSender{}
}

func (f BearyChatSender) IsSupport() bool {
	return utils.Config.SenderConfig.BearyChat.IsEnabled
}

func (f BearyChatSender) Send(notifications []*models.Notification) {
	if !utils.Config.SenderConfig.BearyChat.IsEnabled {
		return
	}
	for _, item := range notifications {
		f.SingleSend(item)
	}
}

func (f BearyChatSender) SingleSend(notification *models.Notification) {
	m := bearychat.Incoming{
		Text:         f.BuildMessage(notification),
		Markdown:     true,
		Notification: "Hello",
	}
	output, _ := m.Build()
	_, err := http.Post(utils.Config.SenderConfig.BearyChat.GroupEndpoint, "application/json", output)
	if err != nil {
		log.Error(err)
	}
}

func (f BearyChatSender) BuildMessage(notification *models.Notification) string {
	return fmt.Sprintf("status:%d,type:%s,url:%s", notification.HttpStatus, notification.Reason, notification.URL)
}
