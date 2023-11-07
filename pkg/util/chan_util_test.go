package util

import (
	"context"
	"testing"
	"time"
)

func TestSendWithBackOffFullChannel(t *testing.T) {
	strCh := make(chan string, 2)
	go func() {
		for s := range strCh {
			time.Sleep(10 * time.Millisecond)
			t.Logf("receive:%s", s)
		}
	}()
	for i := 0; i < 5; i++ {
		go SendWithBackOff(context.Background(), "s", strCh, 2, 10*time.Millisecond)
	}
	time.Sleep(time.Second)
	close(strCh)
}

func TestSendWithBackOff(t *testing.T) {
	strCh := make(chan string, 2)
	go func() {
		for s := range strCh {
			t.Logf("receive:%s", s)
		}
	}()
	for i := 0; i < 5; i++ {
		go SendWithBackOff(context.Background(), "s", strCh, 2, 10*time.Millisecond)
	}
	time.Sleep(time.Second)
	close(strCh)
}
