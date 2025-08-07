package sys_devtools

import (
	"gitee.com/mxcker/go-admin/server/utils"
)

type Filed struct {
	Id          int    `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Type        string `json:"type" binding:"required"`
	ChineseName string `json:"chinese_name" binding:"required"`
	IndexType   string `json:"index_type" binding:"required"`
	TableShow   bool   `json:"table_show" binding:"required"`
	Editable    bool   `json:"editable" binding:"required"`
	Required    bool   `json:"required" binding:"required"`
	Hidden      bool   `json:"hidden" binding:"required"`
}

type GenerateBody struct {
	ModuleName        string  `json:"module_name" binding:"required"`
	ChineseModuleName string  `json:"chinese_module_name" binding:"required"`
	ModelName         string  `json:"model_name" binding:"required"`
	UseCommon         bool    `json:"use_common"`
	UseSoftDelete     bool    `json:"use_soft_delete"`
	CreateCURD        bool    `json:"create_curd"`
	Fields            []Filed `json:"fields" binding:"required"`
}

type WriteItem struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type TemplateItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type SysAutoCode struct {
	utils.DbBaseModel
	ModuleName string `json:"module_name" gorm:"unique"`
	ModelName  string `json:"model_name" gorm:"unique"`
	Form       string `json:"form" json:"form,omitempty"`
}
