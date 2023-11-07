package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
	"webhook/pkg/setting"
)

var (
	MysqlFd *gorm.DB
	err     error
)

type BaseModel struct {
	Id         int       `json:"id" gorm:"primary_key"`
	CreateTime time.Time `json:"create_time"`
	CreateBy   string    `json:"create_by"`
}

func AcquireMySQLDB(mysql setting.MySQLSetting) {
	if MysqlFd, err = gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysql.User,
			mysql.Pass,
			mysql.Address,
			mysql.Db),
	); err != nil {
		log.Fatalf("mysql launch fail:%s", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "" + defaultTableName
	}
	log.Println("MySQL acquire success...")
	MysqlFd.LogMode(false)
	MysqlFd.SingularTable(true)
	MysqlFd.DB().SetMaxIdleConns(mysql.MaxIdleConn)
	MysqlFd.DB().SetMaxOpenConns(mysql.MaxOpenConn)
}

func ReleaseMySQL() {
	if MysqlFd != nil {
		_ = MysqlFd.Close()
	}
}
