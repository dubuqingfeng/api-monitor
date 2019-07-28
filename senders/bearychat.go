package senders

import (
	"fmt"
	"github.com/bearyinnovative/bearychat-go"
	"github.com/dubuqingfeng/api-monitor/models"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// sender
type BearyChatSender struct {
	Sender
	UnSupportType map[string]int
}

// bearychat pusher
var BearyChatPusher BearyChatSender

func init() {
	BearyChatPusher = BearyChatSender{UnSupportType: utils.Config.SenderConfig.BearyChat.UnSupportTypes}
}

// is support
func (f BearyChatSender) IsSupport() bool {
	return utils.Config.SenderConfig.BearyChat.IsEnabled
}

// send
func (f BearyChatSender) Send(notifications []*models.Notification) {
	if !utils.Config.SenderConfig.BearyChat.IsEnabled {
		return
	}
	for _, item := range notifications {
		f.SingleSend(item)
	}
}

// send notification
func (f BearyChatSender) SingleSend(notification *models.Notification) {
	if _, ok := f.UnSupportType[notification.Type]; notification.Type != "" && ok {
		return
	}
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

// build message
func (f BearyChatSender) BuildMessage(notification *models.Notification) string {
	return fmt.Sprintf("status:%d,type:%s,monitor:%s,url:%s", notification.HTTPStatus,
		notification.Reason, utils.Config.MonitorName, notification.URL)
}
