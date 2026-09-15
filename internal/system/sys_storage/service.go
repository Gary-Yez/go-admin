package sys_storage

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/Gary-Yez/go-admin/dberror"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type serviceStruct struct{}

func (*serviceStruct) List(ctx context.Context, req *request.ReqList) (list []SysStorage, total int64, err error) {
	db := req.WithFilter(state.DB().WithContext(ctx).Model(&SysStorage{}), []string{"name", "engine", "enabled"})
	if err = db.Count(&total).Error; err != nil {
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	err = req.WithPagination(req.WithSort(db, []string{"id", "created_at"})).Order("id DESC").Omit("ConfigJSON").Find(&list).Error
	for i := range list {
		list[i].IsDefault = list[i].DefaultSlot != nil
	}
	return
}

// Options 仅返回供上传选择的安全字段，不要求用户拥有存储管理权限。
func Options(ctx context.Context) ([]Option, error) {
	var rows []SysStorage
	if err := state.DB().WithContext(ctx).Select("id", "name", "engine", "default_slot", "enabled").Order("id").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]Option, 0, len(rows))
	for _, row := range rows {
		result = append(result, Option{Id: row.Id, Name: row.Name, Engine: row.Engine, IsDefault: row.DefaultSlot != nil, Enabled: row.Enabled})
	}
	return result, nil
}

func normalize(p Params, engine storage.Engine) (Params, error) {
	for _, v := range []*string{&p.Root, &p.Endpoint, &p.Region, &p.Bucket, &p.AccessKey, &p.PublicURL} {
		*v = strings.TrimSpace(*v)
	}
	switch engine {
	case storage.Local:
		if p.Root == "" {
			return p, errors.New("请填写本地存储目录")
		}
		root, err := filepath.Abs(p.Root)
		if err != nil {
			return p, err
		}
		p = Params{Root: root, PublicURL: p.PublicURL}
	case storage.Tencent:
		p = Params{Endpoint: strings.TrimRight(p.Endpoint, "/"), AccessKey: p.AccessKey, SecretKey: p.SecretKey, PublicURL: p.PublicURL}
	case storage.Aliyun:
		p.Root = ""
		p.PathStyle = false
		p.Endpoint = strings.TrimRight(p.Endpoint, "/")
	case storage.S3:
		p.Root = ""
		p.Endpoint = strings.TrimRight(p.Endpoint, "/")
	default:
		return p, errors.New("不支持的存储引擎")
	}
	store, err := storage.New(p.Config(engine))
	if err != nil {
		return p, err
	}
	if closer, ok := store.(io.Closer); ok {
		_ = closer.Close()
	}
	return p, nil
}
func sameLocation(a, b Params) bool {
	return filepath.Clean(a.Root) == filepath.Clean(b.Root) && strings.TrimRight(a.Endpoint, "/") == strings.TrimRight(b.Endpoint, "/") && a.Region == b.Region && a.Bucket == b.Bucket && a.PathStyle == b.PathStyle
}
func inUse(tx *gorm.DB, id uint) (bool, error) {
	var count int64
	err := tx.Table("sys_files").Where("storage_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (*serviceStruct) Save(ctx context.Context, body *SaveBody) error {
	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" {
		return errors.New("请填写存储名称")
	}
	if body.IsDefault && !body.Enabled {
		return errors.New("默认存储必须启用")
	}
	// 管理操作与未命中回填共用锁；提交前清空缓存，失败则不写数据库。
	lock, err := accountLock(ctx, body.Id)
	if err != nil {
		return err
	}
	defer lock.Unlock()
	if body.Id != 0 {
		if err := state.Cache().Del(cacheKey(body.Id)); err != nil {
			return err
		}
	}
	err = state.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 默认切换先清除原默认，唯一索引保证最多只有一个默认账号。
		if body.IsDefault {
			if err := tx.Model(&SysStorage{}).Where("default_slot = ?", 1).Update("default_slot", nil).Error; err != nil {
				return err
			}
		}
		var row SysStorage
		if body.Id != 0 {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&row, body.Id).Error; err != nil {
				return err
			}
			config, err := row.Config()
			if err != nil {
				return err
			}
			old := ParamsFromConfig(config)
			if body.Engine == row.Engine {
				if body.Params.SecretKey == "" {
					body.Params.SecretKey = old.SecretKey
				}
			}
			params, err := normalize(body.Params, body.Engine)
			if err != nil {
				return err
			}
			used, err := inUse(tx, row.Id)
			if err != nil {
				return err
			}
			if used && (body.Engine != row.Engine || !sameLocation(old, params)) {
				return errors.New("已有文件或上传会话引用此存储，不能更换引擎、桶、区域、地址或目录，请新增存储")
			}
			body.Params = params
		} else {
			var err error
			body.Params, err = normalize(body.Params, body.Engine)
			if err != nil {
				return err
			}
		}
		data, err := json.Marshal(body.Params)
		if err != nil {
			return err
		}
		row.Name = body.Name
		row.Engine = body.Engine
		row.Enabled = body.Enabled
		row.ConfigJSON = string(data)
		row.DefaultSlot = nil
		if body.IsDefault {
			one := uint(1)
			row.DefaultSlot = &one
		}
		if row.Id == 0 {
			return tx.Create(&row).Error
		}
		return tx.Select("Name", "Engine", "Enabled", "ConfigJSON", "DefaultSlot").Updates(&row).Error
	})
	return dberror.Unique(err, &SysStorage{})
}
func (*serviceStruct) Delete(ctx context.Context, id uint) error {
	lock, err := accountLock(ctx, id)
	if err != nil {
		return err
	}
	defer lock.Unlock()
	if err := state.Cache().Del(cacheKey(id)); err != nil {
		return err
	}
	return state.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row SysStorage
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&row, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		used, err := inUse(tx, id)
		if err != nil {
			return err
		}
		if used {
			return errors.New("存储仍被文件或上传会话引用，不能删除")
		}
		return tx.Delete(&row).Error
	})
}

// WithUploadAccount 在短事务中锁定存储并创建文件引用，避免与删除、改位置并发。
// 不在这个事务内传输文件；后续分片通过存储 ID 读取缓存配置。
func WithUploadAccount(ctx context.Context, id uint, fn func(*gorm.DB, *SysStorage, storage.Config) error) error {
	return state.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row SysStorage
		db := tx.Clauses(clause.Locking{Strength: "UPDATE"})
		if id == 0 {
			db = db.Where("default_slot = ?", 1)
		} else {
			db = db.Where("id = ?", id)
		}
		if err := db.First(&row).Error; err != nil {
			return errors.New("请选择有效存储，或先设置默认存储")
		}
		if !row.Enabled {
			return errors.New("该存储已停用，请选择其他存储")
		}
		config, err := row.Config()
		if err != nil {
			return err
		}
		return fn(tx, &row, config)
	})
}

func (*serviceStruct) Get(ctx context.Context, id uint) (*Detail, error) {
	var row SysStorage
	db := state.DB().WithContext(ctx)
	if err := db.First(&row, id).Error; err != nil {
		return nil, err
	}
	var params Params
	if err := json.Unmarshal([]byte(row.ConfigJSON), &params); err != nil {
		return nil, err
	}
	params.SecretKey = ""
	row.IsDefault = row.DefaultSlot != nil
	used, err := inUse(db, row.Id)
	if err != nil {
		return nil, err
	}
	return &Detail{Storage: row, Params: params, InUse: used}, nil
}

// SetState 只更新默认标记或启用状态，不回传、覆盖账号凭据。
func (*serviceStruct) SetState(ctx context.Context, id uint, enabled, setDefault bool) error {
	lock, err := accountLock(ctx, id)
	if err != nil {
		return err
	}
	defer lock.Unlock()
	if err := state.Cache().Del(cacheKey(id)); err != nil {
		return err
	}
	return state.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if setDefault {
			if err := tx.Model(&SysStorage{}).Where("default_slot = ?", 1).Update("default_slot", nil).Error; err != nil {
				return err
			}
		}
		var row SysStorage
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&row, id).Error; err != nil {
			return err
		}
		if setDefault {
			if !row.Enabled {
				return errors.New("请先启用该存储，再设为默认")
			}
			return tx.Model(&row).Update("default_slot", 1).Error
		}
		changes := map[string]any{"enabled": enabled}
		if !enabled {
			changes["default_slot"] = nil
		}
		return tx.Model(&row).Updates(changes).Error
	})
}
