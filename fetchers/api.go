package fetchers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dubuqingfeng/api-monitor/models"
	"github.com/dubuqingfeng/api-monitor/processors"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

// NewAPIFetcher new API fetcher
func NewAPIFetcher() *APIFetcher {
	fetcher := &APIFetcher{
		wg: &sync.WaitGroup{},
		ch: make(chan *models.Process),
		client: &http.Client{
			Timeout: utils.Config.Timeout.Timeout * time.Second,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   utils.Config.Timeout.TLSHandshakeTimeout * time.Second,
				ResponseHeaderTimeout: utils.Config.Timeout.ResponseHeaderTimeout * time.Second,
				ExpectContinueTimeout: utils.Config.Timeout.ExpectContinueTimeout * time.Second,
			},
		},
	}
	return fetcher
}

// API fetcher
type APIFetcher struct {
	wg     *sync.WaitGroup
	client *http.Client
	ch     chan *models.Process
}

// Handle fetch
func (f APIFetcher) Handle() {
	// api monitor
	if utils.Config.APIMonitorEnabled {
		f.fetchMonitoringAPIs()
	}
	// ping api monitor
	if utils.Config.PingAPIMonitorEnabled {
		f.fetchPingAPIs()
	}
	go func() {
		f.wg.Wait()
		close(f.ch)
	}()
	processor := processors.RequestProcessor{}
	for item := range f.ch {
		processor.Process(item)
	}
}

func (f APIFetcher) fetchMonitoringAPIs() {
	// Get all the endpoints
	// TODO support json format.
	endpoints, err := models.GetAllAPIEndpoints()
	if err != nil {
		log.Error(err)
	}
	// Get all the apis
	apis, err := models.GetAllAPIs()
	if err != nil {
		log.Error(err)
	}
	for _, api := range apis {
		var accessEndpointIds map[int64]int
		if api.AccessEndpointIds != "" {
			accessEndpointIds = make(map[int64]int)
			err = json.Unmarshal([]byte(api.AccessEndpointIds), &accessEndpointIds)
			if err != nil {
				log.Error(err)
				continue
			}
		}

		for _, endpoint := range endpoints {
			// api
			stats, ok := accessEndpointIds[endpoint.ID]
			if len(accessEndpointIds) != 0 && ok && stats == utils.LimitAccessEndpointType {
				// limit access
				continue
			}

			if len(accessEndpointIds) != 0 && !ok {
				// limit
				continue
			}
			f.wg.Add(1)
			go f.fetch(endpoint, api)
		}
	}
}

func (f APIFetcher) fetchPingAPIs() {
	pingAPIs, err := models.GetAllPingAPIs()
	if err != nil {
		log.Error(err)
	}
	for _, pingAPI := range pingAPIs {
		f.wg.Add(1)
		go f.fetch(models.APIEndpoint{Endpoint: pingAPI.Endpoint}, pingAPI)
	}
}

// fetch
func (f APIFetcher) fetch(endpoint models.APIEndpoint, api models.API) {
	defer f.wg.Done()
	var buf bytes.Buffer
	buf.WriteString(endpoint.Endpoint)
	buf.WriteString(api.APIURL)
	request, err := http.NewRequest(api.APIMethod, buf.String(), nil)
	if err != nil {
		log.Error(err)
		fmt.Println("Fatal error ", err.Error())
	}

	if api.APIMethod == "GET" && api.QueryStringParams != "" {
		request.URL.RawQuery = api.QueryStringParams
	}

	// request Headers
	if api.RequestHeader != "" {
		headers := make(map[string]string)
		err = json.Unmarshal([]byte(api.RequestHeader), &headers)
		if err != nil {
			log.Error(err)
		}
		for key, value := range headers {
			request.Header.Set(key, value)
		}
	}

	// global environment variable
	globalHeaders := f.GetGlobalENV()
	for key, value := range globalHeaders {
		request.Header.Set(key, value)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("UA", "api-monitor")
	resp, err := f.client.Do(request)
	if err != nil {
		log.Error(err)
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Fatal error ", err.Error())
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()
	process := models.Process{API: api, Endpoint: endpoint, Response: resp, Body: content}
	f.ch <- &process
}

func (f APIFetcher) GetGlobalENV() map[string]string {
	globalHeaders := make(map[string]string)
	globalEnv := os.Getenv("GLOBAL_HEADERS")
	if globalEnv != "" {
		err := json.Unmarshal([]byte(globalEnv), &globalHeaders)
		if err != nil {
			log.Error(err)
		}
	}
	return globalHeaders
}
