package storage

import (
	"fmt"

	"github.com/Gary-Yez/go-admin/storage/drivers/aliyun"
	"github.com/Gary-Yez/go-admin/storage/drivers/local"
	"github.com/Gary-Yez/go-admin/storage/drivers/s3"
	"github.com/Gary-Yez/go-admin/storage/drivers/tencent"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type Engine string

const (
	Local   Engine = "local"
	Tencent Engine = "tencent"
	Aliyun  Engine = "aliyun"
	S3      Engine = "s3"
)

// Config 按 Type 使用对应字段，不读取数据库或缓存，不自动设置默认存储。
// 腾讯云 Endpoint 对应完整 BucketURL，AccessKey 对应 SecretID。
type Config struct {
	Type        Engine
	Root        string // 本地目录。
	Endpoint    string
	Region      string
	Bucket      string
	AccessKey   string
	SecretKey   string
	PublicURL   string
	PathStyle   bool                    // S3 路径式桶地址。
	Credentials aws.CredentialsProvider // 仅 S3：优先于静态凭据。
}

// New 每次创建新适配器。配置读取、存储名称和上传会话由调用方管理。
// 本地适配器用完后通过 io.Closer 关闭；云端连接池复用。
func New(c Config) (Store, error) {
	var instance Store
	var err error
	switch c.Type {
	case Local:
		instance, err = local.New(local.Config{Root: c.Root, PublicURL: c.PublicURL})
	case Tencent:
		instance, err = tencent.New(tencent.Config{BucketURL: c.Endpoint, SecretID: c.AccessKey, SecretKey: c.SecretKey, PublicURL: c.PublicURL})
	case Aliyun:
		instance, err = aliyun.New(aliyun.Config{Endpoint: c.Endpoint, Region: c.Region, Bucket: c.Bucket, AccessKey: c.AccessKey, SecretKey: c.SecretKey, PublicURL: c.PublicURL})
	case S3:
		instance, err = s3.New(s3.Config{Endpoint: c.Endpoint, Region: c.Region, Bucket: c.Bucket, AccessKey: c.AccessKey, SecretKey: c.SecretKey, PublicURL: c.PublicURL, PathStyle: c.PathStyle, Credentials: c.Credentials})
	default:
		return nil, fmt.Errorf("不支持的存储类型：%q", c.Type)
	}
	if err != nil {
		return nil, err
	}
	return instance, nil
}
