package middleware

import (
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterConcurrentAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewRateLimiter(5, time.Minute)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 120; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			switch i % 6 {
			case 0:
				rl.Snapshot()
			case 1:
				rl.Count()
			case 2:
				rl.Reset("10.0.0.1")
			case 3:
				rl.Peek("10.0.0.1")
			case 4:
				rl.ClearAll()
			default:
				rec := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rec)
				c.Request = httptest.NewRequest("GET", "/api/v1/plants", nil)
				c.Request.RemoteAddr = "10.0.0.1:12345"
				rl.Limit()(c)
			}
		}(i)
	}
	close(start)
	wg.Wait()
}
