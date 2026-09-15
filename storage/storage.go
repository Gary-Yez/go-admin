// Package storage 提供统一的存储创建入口和文件操作接口。
package storage

import "github.com/Gary-Yez/go-admin/storage/internal/types"

type (
	Store       = types.Store
	Object      = types.Object
	PutOptions  = types.PutOptions
	URLOptions  = types.URLOptions
	URLProvider = types.URLProvider
)

var (
	ErrNotFound    = types.ErrNotFound
	ErrUnsupported = types.ErrUnsupported
)
