package global

import (
	"fmt"
	"gitee.com/mxcker/go-admin/server/models/configs"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	StartTime = time.Now()
	RootPath  string
	DB        *gorm.DB
	Redis     string
	Config    *configs.Config
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	RootPath = filepath.Join(cwd, "../")
}

func Init() {
	InitConfig()
	InitDB()
}

func InitConfig() {
	viper.SetConfigName("config")
	// 设置配置文件路径（这里设置为当前目录）
	viper.AddConfigPath(".")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return
	}
	var config configs.Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		return
	}
	Config = &config
	fmt.Println("配置获取成功")
}

func InitDB() {
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
		panic("数据库实例创建失败")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("数据库获取失败")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
	DB = db
	fmt.Println("数据库链接成功")
}
