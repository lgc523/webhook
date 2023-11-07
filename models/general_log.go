package models

import (
	"context"
	"github.com/jinzhu/gorm"
	"time"
)

type GeneralLog struct {
	Id         int       `json:"id"`
	CrashTime  time.Time `json:"crash_time"`
	CreateTime time.Time `json:"create_time"`
	Type       string    `json:"type"`
	Pid        string    `json:"pid"`
	SqlText    string    `json:"sql_text"`
}

func (g *GeneralLog) SaveGeneralLog(ctx context.Context, mysql *gorm.DB) error {
	mysql.Save(g)
	return nil
}
