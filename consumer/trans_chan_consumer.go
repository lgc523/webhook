package consumer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	kfk "webhook/middleware/kafka"
	redMod "webhook/middleware/redis"
)

// FlushRedis2Kafka msg from redis chan to kafka
func FlushRedis2Kafka(ctx context.Context, webQueue chan *redMod.DingMsgPayload, k *kfk.Kfk) {
	for {
		select {
		case msg := <-webQueue:
			marshal, _ := json.Marshal(msg)
			message := kafka.Message{Value: marshal}
			if err := k.SendMessage(ctx, message); err != nil {
				log.Printf("TransChanRedis2Kafka.failed:%s", err.Error())
			}
		case <-ctx.Done():
			log.Println("TransChanRedis2Kafka: context canceled, exiting...")
			return
		}
	}
}
