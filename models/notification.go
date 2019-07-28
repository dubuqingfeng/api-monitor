package models

// notification model
type Notification struct {
	HTTPStatus int
	Reason     string
	Type       string
	URL        string
}
