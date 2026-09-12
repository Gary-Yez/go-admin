package config

import (
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"net"
	"net/url"
	"strings"
	"time"
)

type Database struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (c *Database) Normalize() error {
	c.Driver = strings.ToLower(strings.TrimSpace(c.Driver))
	if c.Driver == "" {
		c.Driver = "mysql"
	}
	switch c.Driver {
	case "mysql":
		if c.Port == "" {
			c.Port = "3306"
		}
	case "postgres":
		if c.Port == "" {
			c.Port = "5432"
		}
		if c.SSLMode == "" {
			c.SSLMode = "disable"
		}
	default:
		return fmt.Errorf("不支持的数据库类型 %q，仅支持 mysql 或 postgres", c.Driver)
	}
	if strings.TrimSpace(c.Host) == "" || strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("请配置 database.host 和 database.name")
	}
	return nil
}

func (c Database) DSN() string {
	if c.Driver == "mysql" {
		dsn := mysqlDriver.NewConfig()
		dsn.User, dsn.Passwd = c.Username, c.Password
		dsn.Net, dsn.Addr, dsn.DBName = "tcp", net.JoinHostPort(c.Host, c.Port), c.Name
		dsn.ParseTime, dsn.Loc = true, time.Local
		dsn.Params = map[string]string{"charset": "utf8mb4"}
		return dsn.FormatDSN()
	}
	dsn := url.URL{Scheme: "postgres", User: url.UserPassword(c.Username, c.Password), Host: net.JoinHostPort(c.Host, c.Port), Path: "/" + c.Name}
	query := url.Values{"sslmode": {c.SSLMode}, "TimeZone": {"UTC"}}
	dsn.RawQuery = query.Encode()
	return dsn.String()
}
