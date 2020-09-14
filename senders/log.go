package senders

import (
	"fmt"
	"github.com/dubuqingfeng/api-monitor/models"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
)

// sender
type LogSender struct {
	Sender
}

// log pusher
var LogPusher LogSender

func init() {
	LogPusher = LogSender{}
}

// is support
func (f LogSender) IsSupport() bool {
	return utils.Config.SenderConfig.Log.IsEnabled
}

// send
func (f LogSender) Send(notifications []*models.Notification) {
	if !utils.Config.SenderConfig.Log.IsEnabled {
		return
	}
	for _, item := range notifications {
		f.SingleSend(item)
	}
}

// send notification
func (f LogSender) SingleSend(notification *models.Notification) {
	message := f.BuildMessage(notification)
	log.Info(message)
}

// build message
func (f LogSender) BuildMessage(notification *models.Notification) string {
	return fmt.Sprintf("status:%d,type:%s,monitor:%s,url:%s", notification.HTTPStatus,
		notification.Reason, utils.Config.MonitorName, notification.URL)
}
