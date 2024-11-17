package models

type SysAutoCode struct {
	DbBaseModel
	ModuleName string `json:"module_name" gorm:"unique"`
	ModelName  string `json:"model_name" gorm:"unique"`
	Form       string `json:"form" json:"form,omitempty"`
}
