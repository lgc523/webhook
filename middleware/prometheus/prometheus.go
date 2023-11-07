package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
	"webhook/pkg/util"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "endpoint"},
	)
	LocalCh chan *DurationMsg
)

type DurationMsg struct {
	Method   string
	Path     string
	Duration int
	CreateAt time.Time
}

func TimeDuration(c *gin.Context) {
	// 记录HTTP请求计数
	HttpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath()).Inc()

	// 记录HTTP请求持续时间
	start := time.Now()
	c.Next()
	duration := time.Since(start).Seconds()
	log.Printf("TimeDuration:%s,%s,%f", c.Request.Method, c.FullPath(), duration)
	dMessage := &DurationMsg{
		Method:   c.Request.Method,
		Path:     c.FullPath(),
		Duration: int(duration * 100000),
		CreateAt: time.Now(),
	}
	if c.FullPath() != "/metrics" {
		go util.SendWithBackOff(c, dMessage, LocalCh, 2, 10*time.Second)
	}
	HttpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
}
