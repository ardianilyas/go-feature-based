package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/ardianilyas/go-feature-based/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	clients = make(map[string]*Client)
	mu      sync.Mutex
)

// GetLimiter ambil limiter per-IP
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if c, exists := clients[ip]; exists {
		c.LastSeen = time.Now()
		return c.Limiter
	}

	limiter := rate.NewLimiter(1, 10)
	clients[ip] = &Client{limiter, time.Now()}
	return limiter
}

// CleanupClients hapus client idle
func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.LastSeen) > time.Minute*5 {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// RateLimitMiddleware middleware gin
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			config.Log.WithFields(logrus.Fields{
				"ip": ip,
				"path": c.Request.URL.Path,
				"method": c.Request.Method,
			}).Warn("Rate limit exceeded")
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}