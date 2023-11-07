package consumer

import (
	"regexp"
	"strings"
	"testing"
	"time"
	"webhook/pkg/util"
)

func TestParseGeneralLogMessage(t *testing.T) {
	input := []string{
		"2023-10-09T07:50:52.368998Z\t150551 Query\tCOMMIT",
	}
	for _, s := range input {
		message, err := parseGeneralLogMessage(s)
		if err != nil {
			t.Error(err.Error())
		}
		t.Logf("%+v", message)
	}
}

func TestParse(t *testing.T) {
	logMessage := "2023-10-10T06:04:29.155159Z\t156452 Query\tSELECT\n\t\t  page_size, compress_ops, compress_ops_ok, compress_time, uncompress_ops, uncompress_time\n\t\t  FROM information_schema.innodb_cmp"

	//logMessage := "a\tb c\td"

	// 创建正则表达式模式，匹配制表符或空格分隔的单词，允许"d"为空
	re := regexp.MustCompile(`[\t\s]+|(\s+d)?$`)

	// 使用正则表达式拆分字符串
	parts := re.Split(logMessage, -1)
	timestamp := parts[0]
	_, t2, _ := util.ConvertTime(timestamp, time.RFC3339Nano)
	t.Log(t2)
	t.Log(timestamp)
	threadId := parts[1]
	t.Logf(threadId)
	queryType := parts[2]
	t.Logf(queryType)
	s := parts[3:]
	t.Logf(strings.Join(s, " "))
	//for _, part := range parts {
	//	// 移除可能存在的"d"
	//	part = regexp.MustCompile(`\s+d$`).ReplaceAllString(part, "")
	//	fmt.Println(part)
	//}
}
