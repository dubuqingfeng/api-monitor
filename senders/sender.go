package senders

import "github.com/dubuqingfeng/api-monitor/models"

type Sender interface {
	Send([]*models.Notification)
	IsSupport() bool
}
