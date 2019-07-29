package middleware

import (
	"order-backend/common"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger() gin.HandlerFunc  {
	return LoggerWithWriter()
}

func LoggerWithWriter(notlogged ...string) gin.HandlerFunc  {
	var skip map[string]struct{}

	if length :=len(notlogged);length>0 {
		skip = make(map[string]struct{},length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		//start timer
		start := time.Now()
		path :=c.Request.URL.Path
		raw :=c.Request.URL.RawQuery

		//process request
		c.Next()
		//Log only when path is not being skipped
		if _,ok := skip[path];!ok {
			//stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP:=c.ClientIP()
			method :=c.Request.Method
			statusCode :=c.Writer.Status()
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path =path + "?" +raw
			}

			deviceSerials := c.GetString("device_serials")

			if deviceSerials == "" {
				deviceSerials = "-"
			}
			common.Log.Infof("[GIN] %s | %3d | %13v | %15s | %s %s\n%s",
				deviceSerials,
				statusCode,
				latency,
				clientIP,
				method,
				path,
				comment,
			)

		}
	}
}

