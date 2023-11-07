package consumer

import (
	"context"
	"encoding/json"
	"log"
	"webhook/middleware/ding"
	kfk "webhook/middleware/kafka"
	"webhook/middleware/redis"
	"webhook/models"
)

func KafkaReceive(ctx context.Context, k *kfk.Kfk) {
	for {
		select {
		case <-ctx.Done():
			log.Println("KafkaReceive: context canceled, exiting...")
			return
		default:
			message, err := k.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("receive kafka crash:%s", err)
				return
			}
			msg := &redis.DingMsgPayload{}
			_ = json.Unmarshal(message.Value, msg)
			//log.Printf("KafkaReceive:%s", message.Value)
			_ = models.SaveWebHooks(msg)
			_ = ding.SendDing(msg, ding.DingToken, msg.ProjectKey)
			_ = k.Reader.CommitMessages(ctx, message)
		}
	}
}
