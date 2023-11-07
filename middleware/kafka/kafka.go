package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"log"
	"strings"
	"time"
	"webhook/pkg/setting"
)

type Kfk struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

var (
	KWebhook  *Kfk
	KDuration *Kfk

	KGeneralLog *Kfk
)

func AcquireKafka(topic, group string) *Kfk {
	k := setting.Svc.Kafka
	//mechanism, err := scram.Mechanism(scram.SHA256, k.User, k.Password)
	//if err != nil {
	//	log.Fatalf("AcquireKafka.scram.crash%s", err)
	//}

	//dialer := &kafka.Dialer{
	//	SASLMechanism: mechanism,
	//	TLS:           &tls.Config{},
	//}
	brokerList := strings.Split(k.Address, ",")
	writer := &kafka.Writer{
		Addr: kafka.TCP(brokerList...),
		//Transport: &kafka.Transport{
		//	TLS:  &tls.Config{},
		//	SASL: mechanism,
		//},
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           -1,
		BatchSize:              1,
		MaxAttempts:            2,
		WriteTimeout:           3 * time.Second,
		Compression:            compress.Lz4,
		AllowAutoTopicCreation: false,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				log.Printf("kafka.write.Completion.err:%s", err)
			}
			//else {
			//	log.Printf("kafka.write.Complete.len:%d", len(messages))
			//}
		},
		Async: true,
		Logger: kafka.LoggerFunc(func(format string, args ...interface{}) {
			//log.Printf("Kafka writer: "+format, args...)
		}),
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerList,
		Topic:   topic,
		//Dialer:           dialer,
		MaxWait:          30 * time.Second,
		ReadLagInterval:  -1,
		ReadBatchTimeout: 30 * time.Second,
		StartOffset:      kafka.LastOffset,
		MaxBytes:         1 << 20,
		GroupID:          group,
		//Logger: kafka.LoggerFunc(func(format string, args ...interface{}) {
		//	log.Printf("Kafka read: "+format, args...)
		//}),
	})

	return &Kfk{
		Writer: writer,
		Reader: reader,
	}
}

func (k *Kfk) ReleaseKafka() {
	if k.Writer != nil {
		k.Writer.Close()
	}
	if k.Reader != nil {
		k.Reader.Close()
	}
}

func (k *Kfk) SendMessage(ctx context.Context, message kafka.Message) error {
	if k.Writer == nil {
		return errors.New("kafka writer is not initialized")
	}
	return k.Writer.WriteMessages(ctx, message)
}

func (k *Kfk) ReceiveMessage(ctx context.Context) (kafka.Message, error) {
	if k.Reader == nil {
		return kafka.Message{}, errors.New("kafka reader is not initialized")
	}

	message, err := k.Reader.FetchMessage(ctx)
	if err != nil {
		return kafka.Message{}, errors.New("kafka reader fetch err:" + err.Error())
	}
	//log.Printf("message.poll:\ntopic;%s\npartition:%d\ngroup:%s\noffset:%d\nhw:%d", message.Topic, message.Partition, k.Reader.Config().GroupID, message.Offset, message.HighWaterMark)
	if err = k.Reader.CommitMessages(ctx, message); err != nil {
		return kafka.Message{}, errors.New("kafka reader commit err:" + err.Error())
	}
	return message, nil
}
