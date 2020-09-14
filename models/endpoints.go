package models

import (
	"errors"
	"fmt"
	"github.com/dubuqingfeng/api-monitor/dbs"
	"github.com/dubuqingfeng/api-monitor/utils"
	log "github.com/sirupsen/logrus"
)

// api endpoint model
type APIEndpoint struct {
	ID          int64
	Name        string
	Type        string
	Description string
	Endpoint    string
	ServerID    int64
	CreatedAt   string
	UpdatedAt   string
}

// GetAllAPIEndpoints get all api endpoints, table name: prefix_api_endpoints
func GetAllAPIEndpoints() ([]APIEndpoint, error) {
	return GetAllAPIEndpointsByMySQL()
}

func GetAllAPIEndpointsByMySQL() ([]APIEndpoint, error) {
	conn := "api:config:read"
	var list []APIEndpoint
	if exists := dbs.CheckDBConnExists(conn); !exists {
		return list, errors.New("not found this database." + conn)
	}

	sql := fmt.Sprintf("SELECT id, `name`, `type`, description, endpoint, server_id, created_at, "+
		"updated_at FROM %s", utils.Config.APIConfigDatabaseTablePrefix+"api_endpoints")
	rows, err := dbs.DBMaps[conn].Query(sql)
	if err != nil {
		log.Error(err)
		return list, err
	}
	for rows.Next() {
		var endpoint APIEndpoint
		if err := rows.Scan(&endpoint.ID, &endpoint.Name, &endpoint.Type, &endpoint.Description,
			&endpoint.Endpoint, &endpoint.ServerID, &endpoint.CreatedAt, &endpoint.UpdatedAt); err != nil {
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
