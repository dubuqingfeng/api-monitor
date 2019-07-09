package processors

import (
	"github.com/dubuqingfeng/api-monitor/models"
	"github.com/dubuqingfeng/api-monitor/senders"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type RequestProcessor struct {
}

func (r RequestProcessor) Process(process *models.Process) {
	var notifications []*models.Notification
	url := process.Endpoint.Endpoint + process.API.APIURL
	//url := process.Response.Request.URL.RawQuery

	// http status 500
	if process.Response.StatusCode == http.StatusInternalServerError {
		log.Error(process.Response.Status)
		notification := &models.Notification{HttpStatus: process.Response.StatusCode, Reason: "", URL: url}
		notifications = append(notifications, notification)
	}

	// http status 50x
	if process.Response.StatusCode >= http.StatusInternalServerError {
		log.Error(process.Response.Status)
		notification := &models.Notification{HttpStatus: process.Response.StatusCode, Reason: "", URL: url}
		notifications = append(notifications, notification)
		r.SendNotifications(notifications)
		return
	}

	// timeout > 1s

	// assert
	//notification := &models.Notification{HttpStatus: process.Response.StatusCode, Reason: "", URL: url}
	//notifications = append(notifications, notification)
	// send
	r.SendNotifications(notifications)
}

func (r RequestProcessor) SendNotifications(notifications []*models.Notification) {
	if len(notifications) == 0 {
		return
	}

	// pusher list
	pushers := [2]senders.Sender{senders.BearyChatPusher, senders.SlackPusher}

	var wg sync.WaitGroup
	for _, item := range pushers {
		if item == nil {
			continue
		}
		if !item.IsSupport() {
			continue
		}
		wg.Add(1)
		go func(notifications []*models.Notification) {
			item.Send(notifications)
		}(notifications)
	}
	wg.Wait()
}
