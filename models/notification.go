package models

type Notification struct {
	HttpStatus int
	Reason     string
	Type       string
	URL        string
}
