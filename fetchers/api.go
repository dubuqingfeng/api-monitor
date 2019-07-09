package fetchers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dubuqingfeng/api-monitor/models"
	"github.com/dubuqingfeng/api-monitor/processors"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

func NewAPIFetcher() *APIFetcher {
	fetcher := &APIFetcher{}
	return fetcher
}

type APIFetcher struct {
}

func (f APIFetcher) Handle() {
	// Get all the endpoints
	endpoints, err := models.GetAllAPIEndpoints()
	if err != nil {
		log.Error(err)
	}
	// Get all the apis
	apis, err := models.GetAllAPIs()
	if err != nil {
		log.Error(err)
	}
	ch := make(chan *models.Process)
	var wg sync.WaitGroup
	for _, api := range apis {
		var accessEndpointIds map[int64]int
		if api.AccessEndpointIds != "" {
			accessEndpointIds = make(map[int64]int)
			err = json.Unmarshal([]byte(api.AccessEndpointIds), &accessEndpointIds)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Error(accessEndpointIds)
		}
		for _, endpoint := range endpoints {
			// api
			stats, ok := accessEndpointIds[endpoint.ID]
			if len(accessEndpointIds) != 0 && ok {
				if stats == utils.LimitAccessEndpointType {
					// limit access
					continue
				}
				log.Error(stats)
			}

			if len(accessEndpointIds) != 0 && !ok {
				// limit
				continue
			}
			wg.Add(1)
			go func(endpoint models.APIEndpoint, api models.API, ch chan *models.Process) {
				f.fetch(endpoint, api, ch)
				wg.Done()
			}(endpoint, api, ch)
		}
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	processor := processors.RequestProcessor{}
	for item := range ch {
		log.Info(item)
		processor.Process(item)
	}
}

func (f APIFetcher) fetch(endpoint models.APIEndpoint, api models.API, ch chan *models.Process) error {
	client := &http.Client{Timeout: 10 * time.Second}
	var buf bytes.Buffer
	buf.WriteString(endpoint.Endpoint)
	buf.WriteString(api.APIURL)
	log.Info(buf.String())
	request, err := http.NewRequest(api.APIMethod, buf.String(), nil)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	if api.APIMethod == "GET" && api.QueryStringParams != "" {
		request.URL.RawQuery = api.QueryStringParams
	}

	if api.RequestHeader != "" {
		request.Header.Set("Content-Type", "application/json")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("UA", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info(resp.Status)
	//content, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("Fatal error ", err.Error())
	//}
	//fmt.Println(string(content))
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	process := models.Process{API: api, Endpoint: endpoint, Response: resp}
	ch <- &process
	return nil
}
