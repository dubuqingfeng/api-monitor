package models

import "net/http"

// process
type Process struct {
	Endpoint APIEndpoint
	API      API
	Response *http.Response
	Body     []byte
}
