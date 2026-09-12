package sys_config

import (
	"encoding/json"
	"github.com/Gary-Yez/go-admin/internal/bizconfig"
)

type Field struct {
	bizconfig.Definition
	Name string `json:"name"`
}
type DefinitionFile struct {
	BuiltinGroups []string `json:"builtin_groups"`
	Groups        []string `json:"groups"`
	Fields        []Field  `json:"fields"`
	Hash          string   `json:"hash"`
	typeName      string
}
type SaveBody struct {
	Fields []Field `json:"fields"`
	Hash   string  `json:"hash"`
}

type Preview struct {
	Path   string `json:"path"`
	Before string `json:"before"`
	After  string `json:"after"`
	Action string `json:"action"`
}

type ValueBody struct {
	Key   string          `json:"key" binding:"required"`
	Value json.RawMessage `json:"value"`
}

type CleanupItem struct {
	Key string `json:"key" binding:"required"`
}

type CleanupBody struct {
	Items []CleanupItem `json:"items" binding:"required,min=1,dive"`
}

type SiteInfo struct {
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	Favicon   string `json:"favicon"`
	Title     string `json:"title"`
	Copyright string `json:"copyright"`
}
