package senders

import "github.com/dubuqingfeng/api-monitor/models"

// sender
type Sender interface {
	Send([]*models.Notification)
	IsSupport() bool
}
