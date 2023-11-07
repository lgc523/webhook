package ding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"webhook/middleware/redis"
	"webhook/pkg/util"
)

var DingToken = "cece0598d9fced351fa0df4516c5e273ba2c65560fd85886ff2b1736f49cf711"

// SendDing 发送钉钉
func SendDing(m *redis.DingMsgPayload, dingToken, projectKey string) error {
	log.Printf("SendDing.param:%v\n", m)
	// 成功失败标志
	var picUrl string
	if m.QualityStatus == "SUCCESS" {
		picUrl = "https://bj-test-datachain.oss-cn-beijing.aliyuncs.com/ypsy/syfin/operateData/XWQWkeiZ_success-32.png"
	} else {
		picUrl = "https://bj-test-datachain.oss-cn-beijing.aliyuncs.com/ypsy/syfin/operateData/uQuoP7Gi_fail-32.png"
	}

	msgUrl := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", dingToken)
	messageUrl := fmt.Sprintf("%s/dashboard?id=%s", util.SonarServer, projectKey)

	link := make(map[string]string)
	link["title"] = fmt.Sprintf("%s代码扫描报告push:[%s]", projectKey, m.QualityStatus)
	link["text"] = fmt.Sprintf("Bugs: %s |漏洞: %s |异味: %s |覆盖率: %s%% |重复率: %s%%\n commit %s\n",
		m.Bugs, m.Vulnerabilities, m.CodeSmells, m.Coverage, m.DuplicatedLinesDensity, m.Commit)
	log.Println(link["text"])
	link["messageUrl"] = messageUrl
	link["picUrl"] = picUrl

	param := make(map[string]interface{})
	param["msgtype"] = "link"
	param["link"] = link

	paramBytes, _ := json.Marshal(param)
	dingTalkRsp, _ := http.Post(msgUrl, "application/json", bytes.NewBuffer(paramBytes))
	dingTalkObj := make(map[string]interface{})
	_ = json.NewDecoder(dingTalkRsp.Body).Decode(&dingTalkObj)
	if dingTalkObj["errcode"].(float64) != 0 {
		log.Printf("消息推送失败%s", dingTalkObj["errmsg"])
		return errors.New("消息推送失败:" + dingTalkObj["errmsg"].(string))
	}
	return nil
}
