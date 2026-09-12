package bizconfig

import (
	"encoding/json"
	"time"
)

type Definition struct {
	Sort        int             `json:"-"`
	Group       string          `json:"group"`
	Key         string          `json:"key"`
	Label       string          `json:"label"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Default     json.RawMessage `json:"default"`
}

type SysConfigValue struct {
	Sort         int       `json:"-" gorm:"default:0"`
	Group        string    `json:"group" gorm:"column:config_group;size:100"`
	Key          string    `json:"key" gorm:"column:config_key;size:191;primaryKey"`
	Label        string    `json:"label" gorm:"size:200"`
	Description  string    `json:"description" gorm:"type:text"`
	Type         string    `json:"type" gorm:"size:32"`
	DefaultValue string    `json:"-" gorm:"type:text"`
	Value        string    `json:"-" gorm:"type:text;not null"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (row SysConfigValue) definition() Definition {
	return Definition{Sort: row.Sort, Key: row.Key, Group: row.Group, Label: row.Label, Description: row.Description, Type: row.Type, Default: json.RawMessage(row.DefaultValue)}
}

type ValueRow struct {
	Registered bool `json:"registered"`
	Definition
	Value     json.RawMessage `json:"value"`
	UpdatedAt time.Time       `json:"updated_at"`
}
