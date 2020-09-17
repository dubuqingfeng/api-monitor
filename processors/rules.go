package processors

import (
	"github.com/dubuqingfeng/api-monitor/models"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"time"
)

type RulerLoader struct {
	db *buntdb.DB
}

var Loader RulerLoader

func InitRulerLoader() {
	db, err := buntdb.Open("data.db")
	if err != nil {
		log.Error(err)
	}
	Loader = RulerLoader{db}
}

func NewRulerLoader() RulerLoader {
	return RulerLoader{}
}

func (r RulerLoader) Init() error {
	var err error
	r.db, err = buntdb.Open("data.db")
	if err != nil {
		log.Error(err)
	}
	return err
}

func (r RulerLoader) Close() error {
	return r.db.Close()
}

func (r RulerLoader) HandleNotifications(notifications []*models.Notification, process *models.Process) []*models.Notification {
	var response []*models.Notification
	key := process.Endpoint.Endpoint + process.API.APIURL
	var alert string
	for _, item := range notifications {
		alert = alert + item.Reason
	}
	if len(notifications) == 0 {
		var val string
		var err error
		err = r.db.View(func(tx *buntdb.Tx) error {
			val, err = tx.Get(key)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil && err != buntdb.ErrNotFound {
			log.Error(err)
		}
		if val != "" {
			notification := &models.Notification{HTTPStatus: process.Response.StatusCode, Reason: "恢复正常", URL: key}
			response = append(response, notification)
		}
		alert = ""
	}

	err := r.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, alert, &buntdb.SetOptions{Expires: true, TTL: time.Hour})
		return err
	})
	if err != nil {
		log.Error(err)
	}

	return response
}
