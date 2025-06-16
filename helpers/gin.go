package helpers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetIpAddress 获取请求的IP地址
func GetIpAddress(gctx *gin.Context) string {
	ip := gctx.GetHeader("X-Real-IP")
	if ip == "" {
		ip = gctx.GetHeader("X-Forwarded-For")
	}
	if ip == "" {
		ip = gctx.GetHeader("cf-connecting-ip")
	}
	if ip == "" {
		ip = gctx.ClientIP()
	}
	return ip
}

// 常见的搜索引擎爬虫 User-Agent 关键字
var botUserAgents = []string{
	"Googlebot",
	"Bingbot",
	"Baiduspider",
	"YandexBot",
	"Sogou",
	"360Spider",
	"Yahoo! Slurp",
	"DuckDuckBot",
	"Exabot",
	"Facebot",
	"Twitterbot",
}

// IsBotUserAgent 检查请求是否来自搜索引擎爬虫
func IsBotUserAgent(userAgent string) bool {
	if userAgent == "" {
		return false
	}
	for _, bot := range botUserAgents {
		if strings.Contains(strings.ToLower(userAgent), strings.ToLower(bot)) {
			return true
		}
	}
	return false
}

func IsBotRequest(gctx *gin.Context) (bool, string) {
	userAgent := gctx.GetHeader("User-Agent")
	return IsBotUserAgent(userAgent), userAgent
}
