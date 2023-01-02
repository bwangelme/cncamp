package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var (
	VERSION string
)

func init() {
	VERSION = os.Getenv("VERSION")
}

// Healthz
// 健康检查的接口，返回 ok
func Healthz(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

// Home
// 响应体以 JSON 的方式返回
// VERSION 表示从环境变量中读取的版本号
// Headers 表示所有请求的 header
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"VERSION": VERSION,
		"Headers": c.Request.Header,
	})
}

// LogMiddleware
// 记录请求日志的中间件
func LogMiddleware(c *gin.Context) {
	c.Next()
	status := c.Writer.Status()
	size := c.Writer.Size()
	fmt.Printf("Request Log -- URI: %v IP: %v Status: %v Size: %v\n", c.Request.URL, c.ClientIP(), status, size)
}

func main() {
	r := gin.Default()

	r.Use(LogMiddleware)
	r.GET("/healthz", Healthz)
	r.GET("/", Home)

	r.Run("localhost:8080")
}
