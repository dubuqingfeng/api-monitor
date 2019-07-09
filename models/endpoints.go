package models

import (
	"errors"
	"fmt"
	"github.com/dubuqingfeng/api-monitor/dbs"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
)

type APIEndpoint struct {
	ID          int64
	Name        string
	Type        string
	Description string
	Endpoint    string
	ServerId    int64
	CreatedAt   string
	UpdatedAt   string
}

// get all api endpoints, table name: prefix_api_endpoints
func GetAllAPIEndpoints() ([]APIEndpoint, error) {
	conn := "api:config:read"
	var list []APIEndpoint
	if exists := dbs.CheckDBConnExists(conn); !exists {
		return list, errors.New("not found this database." + conn)
	}

	var sql string
	sql = fmt.Sprintf("SELECT id, `name`, `type`, description, endpoint, server_id, created_at, "+
		"updated_at FROM %s", utils.Config.APIConfigDatabaseTablePrefix+"api_endpoints")
	rows, err := dbs.DBMaps[conn].Query(sql)
	if err != nil {
		log.Error(err)
		return list, err
	}
	for rows.Next() {
		var endpoint APIEndpoint
		if err := rows.Scan(&endpoint.ID, &endpoint.Name, &endpoint.Type, &endpoint.Description,
			&endpoint.Endpoint, &endpoint.ServerId, &endpoint.CreatedAt, &endpoint.UpdatedAt); err != nil {
			log.Error(err)
		}
		list = append(list, endpoint)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return list, err
	}
	return list, nil
}
