package initialization

import (
	"errors"
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	if err := cfg.Database.Normalize(); err != nil {
		return nil, err
	}
	var dialector gorm.Dialector
	switch cfg.Database.Driver {
	case "mysql":
		dialector = mysqlDialector{mysql.Open(cfg.Database.DSN()).(*mysql.Dialector)}
	case "postgres":
		dialector = postgresDialector{postgres.Open(cfg.Database.DSN()).(*postgres.Dialector)}
	default:
		return nil, fmt.Errorf("不支持的数据库类型 %q", cfg.Database.Driver)
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		TranslateError: true,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

// 保留驱动错误中的约束名称，同时兼容 errors.Is 的 GORM 错误判断。
type mysqlDialector struct{ *mysql.Dialector }

func (d mysqlDialector) Translate(err error) error {
	return retainDatabaseError(err, d.Dialector.Translate(err))
}

type postgresDialector struct{ *postgres.Dialector }

func (d postgresDialector) Translate(err error) error {
	return retainDatabaseError(err, d.Dialector.Translate(err))
}

func retainDatabaseError(original, translated error) error {
	if original == translated {
		return original
	}
	return errors.Join(translated, original)
}
