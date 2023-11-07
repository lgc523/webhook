package consumer

import (
	"context"
	"log"
	"webhook/middleware/redis"
	"webhook/pkg/setting"
)

var LocalQueue chan *redis.DingMsgPayload

func FlushChan2Redis(ctx context.Context, ch chan *redis.DingMsgPayload) {
	redisSetting := setting.Svc.Redis
	for {
		select {
		case msg := <-ch:
			_ = redis.Publish(ctx, redisSetting.Topic, msg)
		case <-ctx.Done():
			log.Printf("FlushChan2Redis context canceled, exiting...")
			return
		}
	}
}
