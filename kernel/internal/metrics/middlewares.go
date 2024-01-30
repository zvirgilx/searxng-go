package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {

		// before handle process.
		start := time.Now()
		path := c.Request.URL.Path

		// processing
		c.Next()

		// after process
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if statusCode != http.StatusNotFound {
			ResponseCounter.WithLabelValues(method, path, strconv.Itoa(statusCode)).Observe(time.Since(start).Seconds())
		} else {
			ResponseCounter.WithLabelValues(method, "!others", strconv.Itoa(statusCode)).Observe(time.Since(start).Seconds())
		}
	}
}
