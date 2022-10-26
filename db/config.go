package db

import (
	"encoding/json"
	"fmt"

	"github.com/flanksource/commons/logger"
	v1 "github.com/flanksource/config-db/api/v1"
	"github.com/flanksource/config-db/db/models"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// GetConfigItem returns a single config item result
func GetConfigItem(extType, extID string) (*models.ConfigItem, error) {
	ci := models.ConfigItem{}
	tx := db.Limit(1).Find(&ci, "external_type = ? and external_id  @> ?", extType, pq.StringArray{extID})
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ci, nil
}

// CreateConfigItem inserts a new config item row in the db
func CreateConfigItem(ci *models.ConfigItem) error {
	if err := db.Create(ci).Error; err != nil {
		return err
	}

	return nil
}

// UpdateConfigItem updates all the fields of a given config item row
func UpdateConfigItem(ci *models.ConfigItem) error {
	if err := db.Updates(ci).Error; err != nil {
		return err
	}

	return nil
}

// QueryConfigItems ...
func QueryConfigItems(request v1.QueryRequest) (*v1.QueryResult, error) {
	results := db.Raw(request.Query)
	logger.Tracef(request.Query)
	if results.Error != nil {
		return nil, fmt.Errorf("failed to parse query: %s -> %s", request.Query, results.Error)
	}

	response := v1.QueryResult{
		Results: make([]map[string]interface{}, 0),
	}

	rows, err := results.Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %s -> %s", request.Query, err)
	}

	columns, err := rows.Columns()
	if err != nil {
		logger.Errorf("failed to get column details: %v", err)
	}
	if rows.Next() {
		if err := results.ScanRows(rows, &response.Results); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %s -> %s", request.Query, err)
		}
		for _, col := range columns {
			response.Columns = append(response.Columns, v1.QueryColumn{
				Name: col,
			})
		}
	}

	response.Count = len(response.Results)
	return &response, nil
}

// NewConfigItemFromResult creates a new config item instance from result
func NewConfigItemFromResult(result v1.ScrapeResult) (*models.ConfigItem, error) {
	var dataStr string
	switch data := result.Config.(type) {
	case string:
		dataStr = data
	case []byte:
		dataStr = string(data)
	default:
		bytes, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to marshal: %v", result.Config)
		}
		dataStr = string(bytes)
	}

	ci := &models.ConfigItem{
		ExternalID:   append(result.Aliases, result.ID),
		ID:           result.ID,
		ConfigType:   result.Type,
		ExternalType: &result.ExternalType,
		Account:      &result.Account,
		Region:       &result.Region,
		Zone:         &result.Zone,
		Network:      &result.Network,
		Subnet:       &result.Subnet,
		Name:         &result.Name,
		Source:       &result.Source,
		Tags:         &result.Tags,
		Config:       &dataStr,
	}

	if result.CreatedAt != nil {
		ci.CreatedAt = *result.CreatedAt
	}

	return ci, nil
}

func GetJSON(ci models.ConfigItem) []byte {
	data, err := json.Marshal(ci.Config)
	if err != nil {
		logger.Errorf("Failed to marshal config: %+v", err)
	}
	return data
}
