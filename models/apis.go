package models

import (
	"errors"
	"fmt"
	"github.com/dubuqingfeng/api-monitor/dbs"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
)

// API model
type API struct {
	ID                string
	ParamName         string
	APIId             int64
	RequestHeader     string
	QueryStringParams string
	RequestBody       string
	Extra             string
	Assert            string
	APIURL            string
	APIMethod         string
	APIType           string // default: "json"
	AccessEndpointIds string
	Description       string
	ServerID          int64
	CreatedAt         string
	UpdatedAt         string
}

// GetAllAPIs get all apis
func GetAllAPIs() ([]API, error) {
	conn := "api:config:read"
	var list []API
	if exists := dbs.CheckDBConnExists(conn); !exists {
		return list, errors.New("not found this database." + conn)
	}

	var sql string
	prefix := utils.Config.APIConfigDatabaseTablePrefix
	sql = fmt.Sprintf("select a.id, a.name, a.description, a.request_header, a.query_string_params, "+
		"a.request_body, a.access_endpoint_ids, a.assert, b.url, b.method from %s a "+
		"left join %s b on a.`api_id` = b.`id`;", prefix+"api_params", prefix+"apis")
	rows, err := dbs.DBMaps[conn].Query(sql)
	if err != nil {
		log.Error(err)
		return list, err
	}
	for rows.Next() {
		var api API
		if err := rows.Scan(&api.ID, &api.ParamName, &api.Description, &api.RequestHeader, &api.QueryStringParams,
			&api.RequestBody, &api.AccessEndpointIds, &api.Assert, &api.APIURL, &api.APIMethod); err != nil {
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
