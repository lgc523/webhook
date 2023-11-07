package util

import (
	"log"
	"strconv"
)

func FormatStrF2Int(strDecimal string) (int, error) {
	floatDecimal, err := strconv.ParseFloat(strDecimal, 64)
	if err != nil {
		log.Println("无法解析字符串为浮点数:", err)
		return -1, err
	}
	// 将浮点数乘以100并转换为整数
	intResult := int(floatDecimal * 100)
	return intResult, err
}
