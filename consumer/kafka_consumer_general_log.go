package consumer

import (
	"context"
	"github.com/bytedance/sonic"
	"log"
	"regexp"
	"strings"
	"time"
	kfk "webhook/middleware/kafka"
	"webhook/models"
	"webhook/pkg/util"
)

func KafkaReceiveGeneralLog(ctx context.Context, k *kfk.Kfk) {
	for {
		select {
		case <-ctx.Done():
			log.Println("KafkaReceiveGeneralLog: context canceled, exiting...")
			return
		default:
			message, err := k.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("KafkaReceiveGeneralLog reveive crash:%s", err)
				return
			}
			//log.Printf("KafkaReceiveGeneralLog.message:%s", message.Value)
			m := make(map[string]any)
			err = sonic.Unmarshal(message.Value, &m)
			if err != nil {
				log.Printf("sonic.unmarshal.err:%s", err.Error())
			} else {
				originMsg := m["message"].(string)
				logMessage, err := parseGeneralLogMessage(originMsg)
				if err != nil {
					log.Printf("KafkaReceiveGeneralLog.parseGeneralLogMessage.err:%s", err.Error())
				} else {
					if logMessage != nil {
						_ = logMessage.SaveGeneralLog(ctx, models.MysqlFd)
						if err = k.Reader.CommitMessages(ctx, message); err != nil {
							log.Printf("KafkaReceiveGeneralLog.commit.err:%s", err.Error())
						}
					}
				}
			}
		}
	}
}

func parseGeneralLogMessage(logMessage string) (*models.GeneralLog, error) {
	re := regexp.MustCompile(`[\t\s]+|(\s+d)?$`)

	// 使用正则表达式拆分字符串
	parts := re.Split(logMessage, -1)
	timestamp := parts[0]
	threadId := parts[1]
	queryType := parts[2]
	s := parts[3:]
	var queryString string
	if len(s) > 0 {
		queryString = strings.Join(s, " ")
		if strings.Contains(queryString, "INSERT INTO") {
			if strings.Contains(queryString, "general_log") {
				return nil, nil
			}
		}
	}
	_, t, err := util.ConvertTime(timestamp, time.RFC3339Nano)
	if err != nil {
		log.Printf("parseGeneralLogMessage.err:%s", err.Error())
		return nil, err
	}
	generalLog := &models.GeneralLog{
		CrashTime:  t,
		CreateTime: time.Now(),
		Type:       queryType,
		Pid:        threadId,
		SqlText:    queryString,
	}

	return generalLog, nil
}
