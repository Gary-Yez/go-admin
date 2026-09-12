package sys_devtools

import (
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"time"
)

type Filed struct {
	Id          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Type        string `json:"type" binding:"required"`
	ChineseName string `json:"chinese_name"`
	IndexType   string `json:"index_type"`
	QueryType   string `json:"query_type"`
	Sortable    bool   `json:"sortable"`
	TableShow   bool   `json:"table_show"`
	Editable    bool   `json:"editable"`
	Required    bool   `json:"required"`
}

type GenerateBody struct {
	CreateMenu        bool   `json:"create_menu"`
	MenuName          string `json:"menu_name"`
	MenuParentKey     string `json:"menu_parent_key"`
	MenuIcon          string `json:"menu_icon"`
	menuDefinitions   []sys_menu.Definition
	ModuleName        string  `json:"module_name" binding:"required"`
	ChineseModuleName string  `json:"chinese_module_name" binding:"required"`
	ModelName         string  `json:"model_name" binding:"required"`
	UseSoftDelete     bool    `json:"use_soft_delete"`
	CreateCURD        bool    `json:"create_curd"`
	AllowCreate       bool    `json:"allow_create"`
	AllowEdit         bool    `json:"allow_edit"`
	AllowDelete       bool    `json:"allow_delete"`
	Fields            []Filed `json:"fields" binding:"omitempty,dive"`
	// Explicit overwrite approval: preview path -> SHA-256 of the existing file.
	OverwriteFiles map[string]string `json:"overwrite_files,omitempty"`
}

func (body GenerateBody) MenuDefinitions() []sys_menu.Definition { return body.menuDefinitions }

func (body GenerateBody) QueryFields() []Filed {
	var fields []Filed
	for _, field := range body.Fields {
		if field.QueryType != "" {
			fields = append(fields, field)
		}
	}
	return fields
}

func (body GenerateBody) CanCreate() bool {
	return body.CreateCURD && body.AllowCreate
}

func (body GenerateBody) CanEdit() bool {
	return body.CreateCURD && body.AllowEdit
}

func (body GenerateBody) CanDelete() bool {
	return body.CreateCURD && body.AllowDelete
}

type WriteItem struct {
	Path              string `json:"path"`
	Content           string `json:"content"`
	Action            string `json:"action"`
	ExistingHash      string `json:"existing_hash,omitempty"`
	ExistingContent   string `json:"existing_content,omitempty"`
	RequiresOverwrite bool   `json:"requires_overwrite"`
	registration      bool
}

type TemplateItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type SysAutoCode struct {
	Id         uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt  time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"comment:更新时间"`
	ModuleName string    `json:"module_name" gorm:"unique"`
	ModelName  string    `json:"model_name" gorm:"unique"`
	Form       string    `json:"form,omitempty"`
}

type HistoryQuery struct {
	Page    int    `form:"page" binding:"gte=1"`
	Limit   int    `form:"limit" binding:"gte=1,lte=100"`
	Keyword string `form:"keyword" binding:"max=100"`
}

type DeleteHistoryBody struct {
	Ids          []uint `json:"ids" binding:"required,min=1,max=100,dive,gt=0"`
	DeleteFiles  bool   `json:"delete_files"`
	PreviewToken string `json:"preview_token"`
}

type DeleteHistoryPlan struct {
	Files []WriteItem `json:"files"`
	Token string      `json:"token"`
}

type HistoryItem struct {
	Id                uint      `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	ModuleName        string    `json:"module_name"`
	ModelName         string    `json:"model_name"`
	ChineseModuleName string    `json:"chinese_module_name"`
	FieldCount        int       `json:"field_count"`
	CreateCURD        bool      `json:"create_curd"`
	UseSoftDelete     bool      `json:"use_soft_delete"`
	ConfigValid       bool      `json:"config_valid"`
}

// BuiltinFields 固定生成，不接受客户端修改。
func (body GenerateBody) BuiltinFields() []Filed {
	return []Filed{
		{Name: "id", Key: "Id", Type: "uint", ChineseName: "编号", TableShow: true, Sortable: true},
		{Name: "created_at", Key: "CreatedAt", Type: "time.Time", ChineseName: "创建时间", TableShow: true, Sortable: true},
		{Name: "updated_at", Key: "UpdatedAt", Type: "time.Time", ChineseName: "更新时间", TableShow: true, Sortable: true},
	}
}
