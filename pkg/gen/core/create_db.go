package core

import (
	"fmt"
	"net/url"
	"time"

	"xorm.io/xorm"
)

func CreateTables(engine *xorm.Engine, tables []any) error {
	// 添加需要生成的结构体()
	// tables := []any{
	// 	new(T),
	// }
	// 批量创建表
	for _, table := range tables {
		has, err := engine.IsTableExist(table)
		if err != nil {
			return err
		}
		if !has {
			err := engine.Sync2(table)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type MysqlConf struct {
	User string
	Pwd  string
	Addr string
	Port int
	Db   string
}

func ToDb(conf MysqlConf, tables []any) {
	var engine *xorm.Engine
	var err error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.User, url.QueryEscape(conf.Pwd), conf.Addr, conf.Port, conf.Db)
	fmt.Println(dns)
	engine, err = xorm.NewEngine("mysql", dns)
	if err != nil {
		panic("初始化Mysql失败: " + err.Error())
	}
	engine.SetMaxIdleConns(10)
	engine.SetMaxOpenConns(10)
	engine.SetConnMaxLifetime(60 * time.Second)
	err = CreateTables(engine, tables)
	if err != nil {
		panic("failed to initialize table: " + err.Error())
	}
	_ = engine.Close()
}
