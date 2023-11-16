package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"webhook/consumer"
	"webhook/middleware/redis"
	"webhook/models"
	"webhook/pkg/resp"
	"webhook/pkg/util"
)

const (
	DingToken = "ding_token"
)

var (
	dingToken string
)

func GetWebHooks(ctx *gin.Context) {

	queryParam := make(map[string]string)

	projectKey := ctx.Query("projectKey")
	commit := ctx.Query("commit")

	if projectKey != "" {
		queryParam["projectKey"] = projectKey
	}
	if commit != "" {
		queryParam["commit"] = commit
	}
	list := models.GetWebHooks(ctx, queryParam)
	ctx.JSON(resp.SUCCESS, list)
}

// ReceiveWebhook receive webhook payload
func ReceiveWebhook(ctx *gin.Context) {
	dingToken = ctx.Query(DingToken)
	if dingToken == "" {
		log.Println("sonarQubeWebhook.url.query.dingTalk.absent")
		ctx.JSON(http.StatusBadRequest, "dingTalk absent")
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error reading request body")
		ctx.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(ctx.Request.Body)

	sonarJson := map[string]any{}
	if err = json.Unmarshal(body, &sonarJson); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid request data")
		return
	}

	analysedAt := sonarJson["analysedAt"].(string)
	commitHash := sonarJson["revision"].(string)
	taskId := sonarJson["taskId"].(string)
	alertStatus := sonarJson["status"].(string)

	projectMap := sonarJson["project"]
	qualityGate := sonarJson["qualityGate"]

	var qualityStatus string
	if qualityMap, ok := qualityGate.(map[string]any); ok {
		qualityStatus = qualityMap["status"].(string)
	}

	var projectKey string
	if project, ok := projectMap.(map[string]any); ok {
		projectKey = project["key"].(string)
	}

	// get sonar server
	var metricValueMap map[string]string
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	if metricValueMap, err = util.SonarResult(httpClient, projectKey); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid projectKey data")
		return
	}

	var auditTime string
	if auditTime, _, err = util.ConvertTime(analysedAt, util.Utc8); err != nil {
		log.Printf("convertTime.crash:%s", err)
		ctx.JSON(http.StatusBadRequest, "invalid analysedAt data")
		return
	}

	msg := &redis.DingMsgPayload{
		Bugs:                   metricValueMap["bugs"],
		Ncloc:                  metricValueMap["ncloc"],
		Commit:                 commitHash,
		TaskId:                 taskId,
		Coverage:               metricValueMap["coverage"],
		AuditTime:              auditTime,
		CodeSmells:             metricValueMap["code_smells"],
		AlertStatus:            alertStatus,
		QualityStatus:          qualityStatus,
		Vulnerabilities:        metricValueMap["vulnerabilities"],
		ProjectKey:             projectKey,
		DuplicatedLinesDensity: metricValueMap["duplicated_lines_density"],
	}

	log.Printf("receive.msg:%+v", *msg)
	go util.SendWithBackOff(ctx, msg, consumer.LocalSonarQueue, 2, 5*time.Second)
	ctx.JSON(resp.SUCCESS, nil)
}

// GetWebHook get webhook by id
func GetWebHook(ctx *gin.Context) {
	id := ctx.Query("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("GetWebHook.id illegal:%s", err)
		return
	}
	hook, err := models.GetWebHook(idInt)
	if err != nil {
		ctx.JSON(resp.SUCCESS, err.Error())
		return
	}
	ctx.JSON(resp.SUCCESS, *hook)
}
