package util

import (
	"fmt"
	"time"
)

const Utc8 = "2006-01-02T15:04:05Z0700"

func ConvertTime(ti, format string) (string, time.Time, error) {
	// 解析时间字符串为 UTC 时间
	utcTime, err := time.Parse(format, ti)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return "", time.Time{}, err
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	localTime := utcTime.In(loc)
	auditTime := localTime.Format(time.DateTime)
	return auditTime, localTime, nil
}
