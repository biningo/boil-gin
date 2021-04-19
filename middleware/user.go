package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

/**
*@Author lyer
*@Date 4/19/21 18:22
*@Describe
**/

func UserInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginUserId, _ := strconv.Atoi(c.GetHeader("userId"))
		c.Set("loginUserId", loginUserId)
	}
}
