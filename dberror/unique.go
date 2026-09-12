package dberror

import (
	"errors"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var modelSchemas sync.Map
var constraintLabels sync.Map

// Unique 根据模型的唯一索引和 comment 标签生成字段提示，使用默认 GORM 命名规则。
// 非唯一冲突原样返回；无法识别约束时不向用户暴露数据库错误详情。
func Unique(err error, model any) error {
	if err == nil {
		return nil
	}
	var constraint string
	duplicate := errors.Is(err, gorm.ErrDuplicatedKey)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		duplicate = true
		// MySQL 返回 "... for key 'table.index'"，MariaDB 可能只返回索引名。
		if pos := strings.LastIndex(mysqlErr.Message, " for key "); pos >= 0 {
			constraint = strings.Trim(strings.TrimSpace(mysqlErr.Message[pos+9:]), "'")
			if dot := strings.LastIndexByte(constraint, '.'); dot >= 0 {
				constraint = constraint[dot+1:]
			}
		}
	}
	var postgresErr *pgconn.PgError
	if errors.As(err, &postgresErr) && postgresErr.Code == "23505" {
		duplicate = true
		constraint = postgresErr.ConstraintName
	}
	if !duplicate {
		return err
	}
	if label := uniqueLabel(model, constraint); label != "" {
		return errors.New(label + "已存在")
	}
	return errors.New("数据已存在，请检查唯一字段")
}

// schema.Parse 缓存模型元数据，仅在唯一冲突时读取，不访问数据库。
func uniqueLabel(model any, constraint string) string {
	if constraint == "" || model == nil {
		return ""
	}
	metadata, err := schema.Parse(model, &modelSchemas, schema.NamingStrategy{})
	if err != nil {
		return ""
	}
	// GORM 的索引解析会更新字段元数据，每个模型只执行一次。
	resolve, _ := constraintLabels.LoadOrStore(metadata, sync.OnceValue(func() map[string]string {
		labels := make(map[string]string)
		for _, index := range metadata.ParseIndexes() {
			if index.Class != "UNIQUE" {
				continue
			}
			fields := make([]string, 0, len(index.Fields))
			for _, field := range index.Fields {
				fields = append(fields, fieldLabel(field.Field))
			}
			label := strings.Join(fields, "、")
			if len(fields) > 1 {
				label += "组合"
			}
			labels[index.Name] = label
		}
		for name, unique := range metadata.ParseUniqueConstraints() {
			labels[name] = fieldLabel(unique.Field)
		}
		return labels
	}))
	return resolve.(func() map[string]string)()[constraint]
}

func fieldLabel(field *schema.Field) string {
	if field.Comment != "" {
		return field.Comment
	}
	return field.Name
}
