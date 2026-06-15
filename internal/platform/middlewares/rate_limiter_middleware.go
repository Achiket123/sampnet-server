package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()
		key := "ratelimit:" + ip

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			// Fail open or log error? Standard pattern is log & fail open or return 500.
			// Let's log it and proceed so that Redis downtime doesn't block users.
			c.Next()
			return
		}

		if count == 1 {
			// First request in the window, set the TTL
			_ = redisClient.Expire(ctx, key, window).Err()
		}

		if int(count) > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
