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
	Control     string          `json:"control"`
	Rows        int             `json:"rows"`
	Options     []Option        `json:"options"`
	Default     json.RawMessage `json:"default"`
}

type Option struct {
	Label string          `json:"label"`
	Value json.RawMessage `json:"value"`
}

type SysConfigValue struct {
	Sort         int       `json:"-" gorm:"default:0"`
	Group        string    `json:"group" gorm:"column:config_group;size:100"`
	Key          string    `json:"key" gorm:"column:config_key;size:191;primaryKey"`
	Label        string    `json:"label" gorm:"size:200"`
	Description  string    `json:"description" gorm:"type:text"`
	Type         string    `json:"type" gorm:"size:32"`
	Control      string    `json:"control" gorm:"size:32"`
	Rows         int       `json:"rows" gorm:"default:0"`
	OptionsValue string    `json:"-" gorm:"column:config_options;type:text"`
	DefaultValue string    `json:"-" gorm:"type:text"`
	Value        string    `json:"-" gorm:"type:text;not null"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (row SysConfigValue) definition() Definition {
	options := []Option{}
	if row.OptionsValue != "" {
		_ = json.Unmarshal([]byte(row.OptionsValue), &options)
	}
	return Definition{Sort: row.Sort, Key: row.Key, Group: row.Group, Label: row.Label, Description: row.Description, Type: row.Type, Control: row.Control, Rows: row.Rows, Options: options, Default: json.RawMessage(row.DefaultValue)}
}

type ValueRow struct {
	Registered bool `json:"registered"`
	Definition
	Value     json.RawMessage `json:"value"`
	UpdatedAt time.Time       `json:"updated_at"`
}
