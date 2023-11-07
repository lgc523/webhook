package consumer

import (
	"context"
	redMod "webhook/middleware/redis"
)

var Redis2KQueue chan *redMod.DingMsgPayload

func RedisPubSubConsumer(ctx context.Context, webhookQueue chan *redMod.DingMsgPayload) {
	_ = redMod.Subscribe(ctx, webhookQueue)
}
