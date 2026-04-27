package db

import (
	devopsmodel "devops/model"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"pkg/conf"
	"strings"
	"sync"
	"system/model"
	"time"

	_ "github.com/glebarez/sqlite"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
)

var (
	engine     *xorm.Engine
	engineOnce sync.Once
)

// GetDb 获取数据库引擎单例
func GetDb() *xorm.Engine {
	engineOnce.Do(func() {
		engine = newEngine()
	})
	return engine
}

// newEngine 创建数据库引擎
func newEngine() *xorm.Engine {
	cfg := conf.Db
	var eng *xorm.Engine
	var err error

	switch cfg.UseDb {
	case "sqlite":
		rootDir, err := os.Getwd()
		if err != nil {
			panic("failed to get working directory: " + err.Error())
		}
		dbPath := filepath.Join(rootDir, "db", cfg.Sqlite.Db)
		eng, err = xorm.NewEngine("sqlite", dbPath)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			cfg.Mysql.User,
			url.QueryEscape(cfg.Mysql.Pwd),
			cfg.Mysql.Addr,
			cfg.Mysql.Port,
			cfg.Mysql.Db,
		)
		eng, err = xorm.NewEngine("mysql", dsn)
	case "pgsql":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			cfg.Pgsql.Addr,
			cfg.Pgsql.Port,
			cfg.Pgsql.User,
			cfg.Pgsql.Pwd,
			cfg.Pgsql.Db,
		)
		eng, err = xorm.NewEngine("postgres", dsn)
	default:
		panic("unsupported database type: " + cfg.UseDb)
	}

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 连接池配置
	eng.SetMaxIdleConns(cfg.MaxIdleConn)
	eng.SetMaxOpenConns(cfg.MaxOpenConn)
	connMaxLifetime := cfg.ConnMaxLifetime
	if connMaxLifetime <= 0 {
		connMaxLifetime = 60 // 默认60秒
	}
	eng.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	eng.ShowSQL(cfg.ShowSQL)
	eng.Logger().SetLevel(xlog.LogLevel(cfg.LogLevel))

	// 自动同步表结构
	if err = syncTables(eng); err != nil {
		panic("failed to sync tables: " + err.Error())
	}
	// 初始化数据
	err = InitData(eng)
	if err != nil {
		panic("failed to init data: " + err.Error())
	}

	fmt.Println("✅", cfg.UseDb, "ok!")
	return eng
}

type tableCommenter interface {
	Comment() string
}

// syncTables 同步表结构（幂等操作）
func syncTables(eng *xorm.Engine) error {
	tables := []any{
		new(model.SysUser),
		new(model.SysUserRole),
		new(model.SysUserDept),
		new(model.SysUserPost),
		new(model.SysRole),
		new(model.SysRoleMenu),
		new(model.SysRoleDept),
		new(model.SysDept),
		new(model.SysMenu),
		new(model.SysPost),
		new(devopsmodel.DevProcess),
		new(devopsmodel.DevProcessTag),
		new(devopsmodel.DevProcessTagLink),
		new(devopsmodel.DevProject),
		new(devopsmodel.DevProjectTag),
		new(devopsmodel.DevProjectTagLink),
	}
	if err := eng.Sync2(tables...); err != nil {
		return err
	}
	return syncTableComments(eng, tables)
}

func syncTableComments(eng *xorm.Engine, tables []any) error {
	driverName := eng.DriverName()
	for _, table := range tables {
		tc, ok := table.(tableCommenter)
		if !ok {
			continue
		}
		comment := tc.Comment()
		if comment == "" {
			continue
		}
		safeComment := strings.ReplaceAll(comment, "'", "''")
		tableName := eng.TableName(table)
		var sql string
		switch driverName {
		case "mysql":
			sql = fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", tableName, safeComment)
		case "postgres":
			sql = fmt.Sprintf("COMMENT ON TABLE \"%s\" IS '%s'", tableName, safeComment)
		default:
			continue
		}
		if _, err := eng.Exec(sql); err != nil {
			return fmt.Errorf("failed to set comment on table %s: %w", tableName, err)
		}
	}
	return nil
}

// InitData 初始化数据库数据（仅在表为空时执行）
func InitData(eng *xorm.Engine) error {
	// 检查是否已有数据
	count, err := eng.Count(new(model.SysMenu))
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// 执行初始化 SQL
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	sqlPath := filepath.Join(rootDir, "db", "init.sql")
	if _, err = eng.ImportFile(sqlPath); err != nil {
		return fmt.Errorf("failed to import %s: %w", sqlPath, err)
	}
	log.Println("✅ database initialized")
	return nil
}
