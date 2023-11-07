package models

import (
	"context"
	"time"
	"webhook/middleware/prometheus"
)

type timeDuration struct {
	Model    BaseModel
	Method   string    `json:"method"`
	Path     string    `json:"path"`
	Duration int       `json:"duration"`
	CrashAt  time.Time `json:"crash_at"`
}

func SaveDuration(ctx context.Context, bean *prometheus.DurationMsg) error {
	m := timeDuration{Model: BaseModel{
		CreateBy: "-",
	},
		Method:   bean.Method,
		Path:     bean.Path,
		Duration: bean.Duration,
		CrashAt:  bean.CreateAt,
	}
	MysqlFd.Save(m)
	return nil
}
