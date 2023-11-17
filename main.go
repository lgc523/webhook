package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webhook/consumer"
	"webhook/middleware/kafka"
	prom "webhook/middleware/prometheus"
	"webhook/middleware/redis"
	"webhook/models"
	"webhook/pkg/setting"
	"webhook/routers"
)

func initConfig() {
	setting.LoadEnv()
	setting.LoadServerConfig()

	models.AcquireMySQLDB(setting.Svc.Mysql)
	//redis.AcquireRedis(setting.Svc.Redis)

	//kafkaSetting := setting.Svc.Kafka
	//kafka.KWebhook = kafka.AcquireKafka(kafkaSetting.WebHookTopic, kafkaSetting.WebHookGroup)
	//kafka.KDuration = kafka.AcquireKafka(kafkaSetting.DurationTopic, kafkaSetting.DurationGroup)
	//kafka.KGeneralLog = kafka.AcquireKafka(kafkaSetting.GeneralTopic, kafkaSetting.GeneralGroup)

	prometheus.MustRegister(prom.HttpRequestsTotal)
	prometheus.MustRegister(prom.HttpRequestDuration)
}

var (
	buildTime = ""

	commit = ""

	goVersion = ""
)

func main() {
	args := os.Args
	if len(args) == 2 && (args[1] == "-v" || args[1] == "version") {
		fmt.Printf("Build Time: %s\n", buildTime)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Platform: %s\n", goVersion)
		return
	}
	initConfig()
	handler := routers.InitRouterHandler()
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", setting.Svc.App.Port),
		Handler:        handler,
		MaxHeaderBytes: 1 << setting.Svc.App.MaxHeader,
		ReadTimeout:    setting.Svc.App.ReadTimeOut * time.Second,
		WriteTimeout:   setting.Svc.App.WriteTimeOut * time.Second,
	}

	stopChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())

	webHookConsumer(ctx)
	durationConsumer(ctx)

	//go consumer.KafkaReceiveGeneralLog(ctx, kafka.KGeneralLog)

	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range stopChan {
			cancel()
			release()
		}
	}()
	log.Fatalf("http server launch fail:%s", server.ListenAndServe())
}

func webHookConsumer(ctx context.Context) {
	// local chan 2 redis
	consumer.LocalSonarQueue = make(chan *redis.DingMsgPayload, 20)
	//go consumer.FlushChan2Redis(ctx, consumer.LocalSonarQueue)
	go consumer.FlushChan2MySQL(ctx, consumer.LocalSonarQueue)

	// redis consumer 2 chan
	//consumer.Redis2KQueue = make(chan *redis.DingMsgPayload, 10)
	//go consumer.RedisPubSubConsumer(ctx, consumer.Redis2KQueue)

	// redis 2 kafka
	//go consumer.FlushRedis2Kafka(ctx, consumer.Redis2KQueue, kafka.KWebhook)

	// kafka 2 mysql
	//go consumer.KafkaReceive(ctx, kafka.KWebhook)
}

func durationConsumer(ctx context.Context) {
	prom.LocalDurationCh = make(chan *prom.DurationMsg, 80)
	go consumer.DurationLocalChan2MySQL(ctx, prom.LocalDurationCh)
	//go consumer.DurationLocalChan(ctx, prom.LocalDurationCh, kafka.KDuration)
	//go consumer.LogAppendDuration(ctx, kafka.KDuration)
}

func release() {
	redis.ReleaseRedis()
	kafka.KWebhook.ReleaseKafka()
	kafka.KDuration.ReleaseKafka()
	models.ReleaseMySQL()
}
