package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReleaseKafka(t *testing.T) {
	k := &Kfk{}
	k.ReleaseKafka()
}

func TestKafkaSendReceive(t *testing.T) {
	k := &Kfk{}
	err := k.SendMessage(context.Background(), kafka.Message{Value: []byte("-")})
	assert.Equal(t, "kafka writer is not initialized", err.Error())

	_, err = k.ReceiveMessage(context.Background())
	assert.Equal(t, "kafka reader is not initialized", err.Error())

}
