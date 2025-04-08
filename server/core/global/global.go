package global

import (
	"gitee.com/mxcker/go-admin/server/core/configs"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	Config *configs.Config
	DB     *gorm.DB
	Redis  string
	Vars   variableMap
)

func Init() error {
	err := InitConfig()
	if err != nil {
		return err
	}
	err = InitDB()
	if err != nil {
		return err
	}
	return nil
}

func InitConfig() error {
	viper.SetConfigName("config")
	// 设置配置文件路径（这里设置为当前目录）
	viper.AddConfigPath(".")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	var config configs.Config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	Config = &config
	return nil
}

func InitDB() error {
	// 配置数据库
	db, err := gorm.Open(mysql.Open(Config.Mysql.ToString()), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
	DB = db
	return nil
}
