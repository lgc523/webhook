package models

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	mock sqlmock.Sqlmock
	db   *sql.DB
	gdb  *gorm.DB
)

func TestModel(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {

		panic(err)
	}
	gdb, err = gorm.Open("mysql", "db")
}
