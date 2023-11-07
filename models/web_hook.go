package models

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
	"webhook/middleware/redis"
	"webhook/pkg/util"
)

type WebHook struct {
	BaseModel
	Bugs            int    `json:"bugs"`
	Vulnerabilities int    `json:"vulnerabilities"`
	BadSmell        int    `json:"bad_smell"`
	Coverage        int    `json:"coverage"`
	Duplicate       int    `json:"duplicate"`
	Status          string `json:"status"`
	QualityStatus   string `json:"quality_status"`
	ProjectKey      string `json:"projectKey"`
	Hash            string `json:"hash"`
}

func GetWebHooks(ctx *gin.Context, condition map[string]string) (list []WebHook) {
	offset, limit := util.GetPage(ctx)
	query := MysqlFd
	if len(condition) > 0 {
		query = query.Where(condition)
	}
	query.Offset(offset).Limit(limit).Find(&list)
	return
}

func SaveWebHooks(msg *redis.DingMsgPayload) error {
	bugs, _ := strconv.Atoi(msg.Bugs)
	cov, _ := util.FormatStrF2Int(msg.Coverage)
	dup, _ := util.FormatStrF2Int(msg.DuplicatedLinesDensity)
	vulnerabilities, _ := util.FormatStrF2Int(msg.Vulnerabilities)
	codeSmell, _ := strconv.Atoi(msg.CodeSmells)
	if msg.AlertStatus == "SUCCESS" {
		msg.AlertStatus = "Y"
	} else {
		msg.AlertStatus = "N"
	}
	if msg.QualityStatus == "SUCCESS" {
		msg.QualityStatus = "Y"
	} else {
		msg.QualityStatus = "N"
	}
	hook := &WebHook{BaseModel: BaseModel{
		CreateBy:   msg.Commit,
		CreateTime: time.Now(),
	},
		Bugs:            bugs,
		Vulnerabilities: vulnerabilities,
		BadSmell:        codeSmell,
		Coverage:        cov,
		Duplicate:       dup,
		Status:          msg.AlertStatus,
		QualityStatus:   msg.QualityStatus,
		ProjectKey:      msg.ProjectKey,
		Hash:            msg.Commit,
	}
	save := MysqlFd.Save(hook)
	log.Infof("SaveWebhook.result:%d", save.RowsAffected)
	return nil
}

func GetWebHook(id int) (*WebHook, error) {
	var wh WebHook
	result := MysqlFd.First(&wh, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("webhook absent")
		}
		if errors.Is(result.Error, gorm.ErrInvalidSQL) {
			return nil, errors.New("illegal sql")
		}
		return nil, result.Error
	}

	return &wh, nil
}
