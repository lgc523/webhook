package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"webhook/pkg/setting"
)

// WebhookPayload defines the structure of the data expected
// to be received from Redis, including URL, Webhook ID, and relevant data.

type DingMsgPayload struct {
	// alert_status
	// Vulnerabilities
	// bugs
	// code_smells
	// coverage
	// duplicated_lines_density
	// ncloc
	Bugs                   string
	Ncloc                  string
	Commit                 string
	AuditId                string
	TaskId                 string
	Coverage               string
	AuditTime              string
	CodeSmells             string
	AlertStatus            string
	QualityStatus          string
	Vulnerabilities        string
	DuplicatedLinesDensity string
	ProjectKey             string
}

var Red *redis.Client

func AcquireRedis(r setting.RedisSetting) {
	Red = redis.NewClient(
		&redis.Options{
			Addr:     r.Address,
			DB:       r.DB,
			Password: r.Pass,
		})
	pingCmd := Red.Ping(context.Background())
	if err := pingCmd.Err(); err != nil {
		log.Fatalf("AcquireRedis.fail:%s", err)
	}
}
func Subscribe(ctx context.Context, ch chan *DingMsgPayload) error {
	// Subscribe to the "webhooks" channel in Redis
	redisSetting := setting.Svc.Redis
	pubSub := Red.Subscribe(ctx, redisSetting.Topic)

	// Ensure that the PubSub connection is closed when the function exits
	defer func() {
		if err := pubSub.Close(); err != nil {
			log.Println("Error closing PubSub:", err)
		}
	}()
	// Infinite loop to continuously receive messages from the "webhooks" channel
	for {
		select {
		case <-ctx.Done():
			log.Println("Subscribe: context canceled, exiting...")
			return nil
		default:
			// Receive a message from the channel
			msg, err := pubSub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error Subscribe.receiveMessage:%v", err)
				return err // Return the error if there's an issue receiving the message
			}

			payload := &DingMsgPayload{}
			// Unmarshal the JSON payload into the WebhookPayload structure
			err = json.Unmarshal([]byte(msg.Payload), payload)
			if err != nil {
				log.Printf("Error unmarshalling payload:%v", err)
				continue // Continue with the next message if there's an error unmarshalling
			}
			ch <- payload // Sending the payload to the channel
		}
	}
}

func ReleaseRedis() {
	if Red != nil {
		err := Red.Close()
		if err != nil {
			return
		}
	}
}

func Publish(ctx context.Context, channel string, msg *DingMsgPayload) error {
	//log.Printf("Publish.msg:%+v", *msg)
	serializedMsg, err := json.Marshal(*msg)
	if err != nil {
		log.Printf("JSON marshaling error: %v", err)
		return err
	}
	if err = Red.Publish(ctx, channel, serializedMsg).Err(); err != nil {
		log.Printf("redis.Publish.err:%sin%s", err, channel)
		return err
	}
	return nil
}
