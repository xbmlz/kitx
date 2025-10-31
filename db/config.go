package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	oracle "github.com/godoes/gorm-oracle"
)

type Config struct {
	Dialect  string `mapstructure:"dialect" yaml:"dialect" json:"dialect"`
	Host     string `mapstructure:"host" yaml:"host" json:"host"`
	Port     int    `mapstructure:"port" yaml:"port" json:"port"`
	Username string `mapstructure:"username" yaml:"username" json:"username"`
	Password string `mapstructure:"password" yaml:"password" json:"password"`
	Database string `mapstructure:"database" yaml:"database" json:"database"`
	Params   string `mapstructure:"params" yaml:"params" json:"params"`
	ShowSQL  bool   `mapstructure:"show_sql" yaml:"show_sql" json:"show_sql"`
}

func (c *Config) DSN() string {
	switch c.Dialect {
	case "sqlite3":
		// https://github.com/glebarez/sqlite
		dir := filepath.Dir(c.Database)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
		url := c.Database
		if c.Params != "" {
			url += "?" + c.Params
		}
		return url
	case "mysql":
		// https://github.com/go-sql-driver/mysql
		url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
		if c.Params != "" {
			url += "?" + c.Params
		}
		return url
	case "postgres":
		// https://github.com/jackc/pgx
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s", c.Host, c.Port, c.Username, c.Password, c.Database, c.Params)
	case "mssql":
		// https://github.com/microsoft/go-mssqldb
		url := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", c.Username, c.Password, c.Host, c.Port, c.Database)
		if c.Params != "" {
			url += fmt.Sprintf("&%s", c.Params)
		}
		return url
	case "oracle":
		// https://github.com/godoes/gorm-oracle
		options := map[string]string{
			"CONNECTION TIMEOUT": "90",
			"LANGUAGE":           "SIMPLIFIED CHINESE",
			"TERRITORY":          "CHINA",
			"SSL":                "false",
		}
		if c.Params != "" {
			// params = xxx=xxx;yyy=yyy
			params := strings.SplitSeq(c.Params, ";")
			for p := range params {
				kv := strings.Split(p, "=")
				if len(kv) == 2 {
					options[kv[0]] = kv[1]
				}
			}
		}
		return oracle.BuildUrl(c.Host, c.Port, c.Database, c.Username, c.Password, options)
	default:
		panic("database dialect not supported, supported: sqlite3, mysql, postgres, mssql, oracle")
	}
}
