package util

import (
	"context"
	"log"
	"time"
)

func SendWithBackOff[T any](c context.Context, msg T, ch chan T, maxRetries int, retryDelay time.Duration) {

	retries := 0
	timer := time.NewTimer(retryDelay)
	for {
		select {
		case ch <- msg:
			log.Printf("SendWithBackOff.send2Ch over,len:%d", len(ch))
			return
		case <-c.Done():
			return
		case <-timer.C:
			log.Printf("Retrying to send to channel, len:%d", len(ch))
			if retries >= maxRetries {
				// Trigger an alert (e.g., ding alert) here
				log.Println("max retries reached, unable to insert")
				return
			}
			timer.Reset(retryDelay)
			retryDelay *= 2
			retries++
		}
	}
}
