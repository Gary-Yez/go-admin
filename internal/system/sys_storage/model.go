package sys_storage

import (
	"encoding/json"
	"time"

	"github.com/Gary-Yez/go-admin/storage"
)

type SysStorage struct {
	Id        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Name      string         `json:"name" gorm:"size:191;uniqueIndex;comment:存储名称"`
	Engine    storage.Engine `json:"engine" gorm:"size:32;comment:存储引擎"`
	Enabled   bool           `json:"enabled" gorm:"comment:允许新上传"`
	// 唯一索引允许多个 NULL，只有默认存储保存 1，兼容 MySQL 和 PostgreSQL。
	DefaultSlot *uint  `json:"-" gorm:"uniqueIndex;comment:默认存储"`
	ConfigJSON  string `json:"-" gorm:"type:text"`
	IsDefault   bool   `json:"is_default" gorm:"-"`
}

// Params 不暴露 SDK 的运行时对象，支持数据库持久化。
type Params struct {
	Root      string `json:"root"`
	Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	PublicURL string `json:"public_url"`
	PathStyle bool   `json:"path_style"`
}

func (p Params) Config(engine storage.Engine) storage.Config {
	return storage.Config{Type: engine, Root: p.Root, Endpoint: p.Endpoint, Region: p.Region, Bucket: p.Bucket, AccessKey: p.AccessKey,
		SecretKey: p.SecretKey, PublicURL: p.PublicURL, PathStyle: p.PathStyle}
}
func ParamsFromConfig(c storage.Config) Params {
	return Params{Root: c.Root, Endpoint: c.Endpoint, Region: c.Region, Bucket: c.Bucket, AccessKey: c.AccessKey,
		SecretKey: c.SecretKey, PublicURL: c.PublicURL, PathStyle: c.PathStyle}
}

type SaveBody struct {
	Id        uint           `json:"id"`
	Name      string         `json:"name" binding:"required,max=191"`
	Engine    storage.Engine `json:"engine" binding:"required,oneof=local tencent aliyun s3"`
	Enabled   bool           `json:"enabled"`
	IsDefault bool           `json:"is_default"`
	Params    Params         `json:"params"`
}
type Option struct {
	Enabled   bool           `json:"enabled"`
	Id        uint           `json:"id"`
	Name      string         `json:"name"`
	Engine    storage.Engine `json:"engine"`
	IsDefault bool           `json:"is_default"`
}

type Detail struct {
	Storage SysStorage `json:"storage"`
	Params  Params     `json:"params"`
	InUse   bool       `json:"in_use"`
}

func (row *SysStorage) Config() (storage.Config, error) {
	var params Params
	if err := json.Unmarshal([]byte(row.ConfigJSON), &params); err != nil {
		return storage.Config{}, err
	}
	return params.Config(row.Engine), nil
}
