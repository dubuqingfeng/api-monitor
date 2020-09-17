package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dubuqingfeng/api-monitor/dbs"
	"github.com/dubuqingfeng/api-monitor/utils"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

// get all ping apis.
func GetAllPingAPIs() ([]API, error) {
	if utils.Config.SourceFormat.Ping == "json" {
		return GetAllPingAPIsByJSON()
	} else if utils.Config.SourceFormat.Ping == "mysql" {
		return GetAllPingAPIsByMySQL()
	} else if utils.Config.SourceFormat.Ping == "csv" {
		return GetAllPingAPIsByCSV()
	} else {
		return GetAllPingAPIsByMySQL()
	}
}

func GetAllPingAPIsByCSV() ([]API, error) {
	var list []API
	content, err := ioutil.ReadFile("configs/ping.csv")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if utils.HasBOM(content) {
		content = utils.StripBOM(content)
	}
	// var records []*Record
	if err := gocsv.UnmarshalBytes(content, &list); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return list, nil
}

func GetAllPingAPIsByJSON() ([]API, error) {
	file, _ := ioutil.ReadFile("configs/ping.json")
	var list []API
	if err := json.Unmarshal([]byte(file), &list); err != nil {
		log.Error(err)
		return list, err
	}
	return list, nil
}

func GetAllPingAPIsByMySQL() ([]API, error) {
	conn := "api:config:read"
	var list []API
	if exists := dbs.CheckDBConnExists(conn); !exists {
		return list, errors.New("not found this database." + conn)
	}

	var sql string
	prefix := utils.Config.APIConfigDatabaseTablePrefix
	sql = fmt.Sprintf("select a.id, a.name, a.description, a.request_header, a.query_string_params, "+
		"a.request_body, a.access_endpoint_ids, a.assert, a.url, a.method, a.endpoint from %s a;", prefix+"api_ping_params")
	rows, err := dbs.DBMaps[conn].Query(sql)
	if err != nil {
		log.Error(err)
		return list, err
	}
	for rows.Next() {
		var api API
		if err := rows.Scan(&api.ID, &api.ParamName, &api.Description, &api.RequestHeader, &api.QueryStringParams,
			&api.RequestBody, &api.AccessEndpointIds, &api.Assert, &api.APIURL, &api.APIMethod, &api.Endpoint); err != nil {
			log.Error(err)
		}
		list = append(list, api)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return list, err
	}
	return list, nil
}
