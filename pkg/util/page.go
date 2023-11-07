package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetPage(ctx *gin.Context) (offset, limit int) {
	defaultOffset := 0
	defaultLimit := 10

	page := ctx.DefaultQuery("page", "1")
	size := ctx.DefaultQuery("pageSize", strconv.Itoa(defaultLimit))

	pInt, _ := strconv.Atoi(page)
	limit, _ = strconv.Atoi(size)

	if pInt > 0 {
		offset = (pInt - 1) * limit
	} else {
		offset = defaultOffset
	}

	return offset, limit
}
