package models

import "net/http"

type Process struct {
	Endpoint APIEndpoint
	API      API
	Response *http.Response
	Body     []byte
}
