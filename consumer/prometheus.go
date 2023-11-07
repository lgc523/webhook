package consumer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	kk "webhook/middleware/kafka"
	prom "webhook/middleware/prometheus"
	"webhook/models"
)

func DurationLocalChan(ctx context.Context, localCh chan *prom.DurationMsg, k *kk.Kfk) {
	for {
		select {
		case dMessage := <-localCh:
			log.Printf("DurationLocalChan.receive.msg:%+v", dMessage)
			data, _ := json.Marshal(dMessage)
			err := k.SendMessage(ctx, kafka.Message{Value: data})
			if err != nil {
				log.Printf("TimeDuration.kafka.write.err:%s", err)
			}
		case <-ctx.Done():
			log.Println("DurationLocalChan: context canceled, exiting...")
			return
		}
	}

}

func LogAppendDuration(ctx context.Context, kfk *kk.Kfk) {
	for {
		select {
		case <-ctx.Done():
			log.Println("LogAppendDuration: context canceled, exiting...")
			return
		default:
			message, err := kfk.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("LogAppendDuration.receive.err:%s", err)
				return
			}
			value := message.Value
			timeDuration := &prom.DurationMsg{}
			_ = json.Unmarshal(value, timeDuration)
			err = models.SaveDuration(ctx, timeDuration)
			if err != nil {
				log.Printf("LogAppendDuration.saveModel.err:%s", err)
				return
			}
			_ = kfk.Reader.CommitMessages(ctx, message)
		}
	}
}
