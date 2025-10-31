package db

import (
	"log"
	"maps"
	"os"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	dbs  sync.Map
	once sync.Once
	cfg  = make(map[string]Config)
)

func Register(c map[string]Config) {
	once.Do(func() {
		maps.Copy(cfg, c)
	})
}

func Get(key string) *gorm.DB {
	if db, ok := dbs.Load(key); ok {
		return db.(*gorm.DB)
	}
	cfg, ok := cfg[key]
	if !ok {
		panic("database config not found")
	}
	db, err := open(cfg)
	if err != nil {
		panic(err)
	}
	dbs.Store(key, db)
	return db
}

func open(cfg Config) (*gorm.DB, error) {
	dsn := cfg.DSN()
	var dialector gorm.Dialector
	switch cfg.Dialect {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite3":
		dialector = sqlite.Open(dsn)
	case "mssql":
		dialector = sqlserver.Open(dsn)
	case "oracle":
		dialector = oracle.New(oracle.Config{
			DSN:                     dsn,
			IgnoreCase:              false, // query conditions are not case-sensitive
			NamingCaseSensitive:     true,  // whether naming is case-sensitive
			VarcharSizeIsCharLength: true,  // whether VARCHAR type size is character length, defaulting to byte length

			// RowNumberAliasForOracle11 is the alias for ROW_NUMBER() in Oracle 11g, defaulting to ROW_NUM
			RowNumberAliasForOracle11: "ROW_NUM",
		})
	default:
		panic("database driver not supported, supported: sqlite3, mysql, postgres, mssql, oracle")
	}
	logLevel := logger.Silent
	if cfg.ShowSQL {
		logLevel = logger.Info
	}
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	dbConfig := &gorm.Config{
		Logger: dbLogger,
	}
	if cfg.Dialect == "oracle" {
		// 是否禁用默认在事务中执行单次创建、更新、删除操作
		dbConfig.SkipDefaultTransaction = true
		// 是否禁止在自动迁移或创建表时自动创建外键约束
		dbConfig.DisableForeignKeyConstraintWhenMigrating = true
		// 自定义命名策略
		dbConfig.NamingStrategy = schema.NamingStrategy{
			NoLowerCase:         true, // 是否不自动转换小写表名
			IdentifierMaxLength: 30,   // Oracle: 30, PostgreSQL:63, MySQL: 64, SQL Server、SQLite、DM: 128
		}
		// 创建并缓存预编译语句，启用后可能会报 ORA-01002 错误
		dbConfig.PrepareStmt = false
		// 插入数据默认批处理大小
		dbConfig.CreateBatchSize = 50
	}
	db, err := gorm.Open(dialector, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	if cfg.Dialect == "oracle" {
		_, _ = oracle.AddSessionParams(sqlDB, map[string]string{
			"TIME_ZONE":               "+08:00",                       // ALTER SESSION SET TIME_ZONE = '+08:00';
			"NLS_DATE_FORMAT":         "YYYY-MM-DD",                   // ALTER SESSION SET NLS_DATE_FORMAT = 'YYYY-MM-DD';
			"NLS_TIME_FORMAT":         "HH24:MI:SSXFF",                // ALTER SESSION SET NLS_TIME_FORMAT = 'HH24:MI:SS.FF3';
			"NLS_TIMESTAMP_FORMAT":    "YYYY-MM-DD HH24:MI:SSXFF",     // ALTER SESSION SET NLS_TIMESTAMP_FORMAT = 'YYYY-MM-DD HH24:MI:SS.FF3';
			"NLS_TIME_TZ_FORMAT":      "HH24:MI:SS.FF TZR",            // ALTER SESSION SET NLS_TIME_TZ_FORMAT = 'HH24:MI:SS.FF3 TZR';
			"NLS_TIMESTAMP_TZ_FORMAT": "YYYY-MM-DD HH24:MI:SSXFF TZR", // ALTER SESSION SET NLS_TIMESTAMP_TZ_FORMAT = 'YYYY-MM-DD HH24:MI:SS.FF3 TZR';
		})
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(10 * time.Second)
	return db, nil
}
