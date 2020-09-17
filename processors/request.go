package processors

import (
	"encoding/json"
	"github.com/dubuqingfeng/api-monitor/models"
	"github.com/dubuqingfeng/api-monitor/pkg/jsonpath"
	"github.com/dubuqingfeng/api-monitor/senders"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sync"
)

// processor
type RequestProcessor struct {
}

// process
func (r RequestProcessor) Process(process *models.Process) {
	var notifications []*models.Notification
	url := process.Endpoint.Endpoint + process.API.APIURL
	// http status 500
	if process.Response.StatusCode == http.StatusInternalServerError {
		log.Error(process.Response.Status)
		notification := &models.Notification{HTTPStatus: process.Response.StatusCode, Reason: "", URL: url}
		notifications = append(notifications, notification)
	}

	// http status 50x
	if process.Response.StatusCode >= http.StatusInternalServerError {
		log.Error(process.Response.Status)
		notification := &models.Notification{HTTPStatus: process.Response.StatusCode, Reason: "", URL: url}
		notifications = append(notifications, notification)
		r.SendNotifications(notifications)
		return
	}

	// timeout > 1s

	// assert
	assertNotifications := r.ProcessAssert(process)
	notifications = append(notifications, assertNotifications...)
	// send
	notifications = Loader.HandleNotifications(notifications, process)
	r.SendNotifications(notifications)
}

// process assert
func (r RequestProcessor) ProcessAssert(process *models.Process) []*models.Notification {
	var assert models.Assert
	var notifications []*models.Notification
	url := process.Endpoint.Endpoint + process.API.APIURL
	if process.API.Assert == "" {
		return notifications
	}
	err := json.Unmarshal([]byte(process.API.Assert), &assert)
	if err != nil {
		log.Error(err)
	}
	// http response status
	if len(assert.Status) != 0 {
		for _, assertStatus := range assert.Status {
			true := r.ProcessStatusAssert(process, assertStatus)
			if !true {
				notification := &models.Notification{HTTPStatus: process.Response.StatusCode, Reason: "", URL: url,
					Type: utils.APITypeAssertFailed}
				notifications = append(notifications, notification)
			}
		}
	}
	// http response body
	if len(assert.Body) != 0 {
		for _, assertBody := range assert.Body {
			true := r.ProcessBodyAssert(process, assertBody)
			if !true {
				notification := &models.Notification{HTTPStatus: process.Response.StatusCode, Reason: "", URL: url,
					Type: utils.APITypeAssertFailed}
				notifications = append(notifications, notification)
			}
		}
	}
	// http response json path
	if len(assert.JSONPath) != 0 {
		for _, assertStatus := range assert.JSONPath {
			true := r.ProcessJsonPathAssert(process, assertStatus)
			if !true {
				notification := &models.Notification{HTTPStatus: process.Response.StatusCode, Reason: "", URL: url,
					Type: utils.APITypeAssertFailed}
				notifications = append(notifications, notification)
			}
		}
	}
	// http response headers
	// http response cookies
	return notifications
}

// process status assert
func (r RequestProcessor) ProcessStatusAssert(process *models.Process, assert models.AssertItem) bool {
	if assert.Type == "equals" && strconv.Itoa(process.Response.StatusCode) != assert.Value {
		return false
	}
	return true
}

// process body assert
func (r RequestProcessor) ProcessBodyAssert(process *models.Process, assert models.AssertItem) bool {
	if assert.Type == "equals" && string(process.Body) != assert.Value {
		return false
	}
	return true
}

// json path assert
func (r RequestProcessor) ProcessJsonPathAssert(process *models.Process, assert models.AssertItem) bool {
	if assert.Key == "" {
		return true
	}
	if assert.Value == "" {
		return true
	}
	if assert.Type == "equals" {
		value := utils.CastType(assert.Value, assert.ValueType)
		boolean, err := jsonpath.Equal(process.Body, assert.Key, value)
		if err != nil {
			log.Error(err)
		}
		return boolean
	}
	if assert.Type == "not_equals" {
		value := utils.CastType(assert.Value, assert.ValueType)
		boolean, err := jsonpath.NotEqual(process.Body, assert.Key, value)
		if err != nil {
			log.Error(err)
		}
		return boolean
	}
	if assert.Type == "contains" {
		value := utils.CastType(assert.Value, assert.ValueType)
		boolean, err := jsonpath.Contains(process.Body, assert.Key, value)
		if err != nil {
			log.Error(err)
		}
		return boolean
	}
	if assert.Type == "len" {
		value := utils.CastType(assert.Value, "int").(int)
		boolean, err := jsonpath.Len(process.Body, assert.Key, value)
		if err != nil {
			log.Error(err)
		}
		return boolean
	}
	return true
}

// send notifications
func (r RequestProcessor) SendNotifications(notifications []*models.Notification) {
	if len(notifications) == 0 {
		return
	}

	// pusher list
	pushers := [3]senders.Sender{senders.BearyChatPusher, senders.SlackPusher, senders.LogPusher}

	var wg sync.WaitGroup
	for _, item := range pushers {
		if item == nil {
			continue
		}
		if !item.IsSupport() {
			continue
		}
		wg.Add(1)
		go func(notifications []*models.Notification, sender senders.Sender) {
			sender.Send(notifications)
			wg.Done()
		}(notifications, item)
	}
	wg.Wait()
}
