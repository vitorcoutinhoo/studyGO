package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

type Client struct {
	Requests int
	ResetAt  time.Time
}

type RateLimiter struct {
	clients map[string]*Client
	mu      sync.Mutex
	limit   int
	window  time.Duration
}

type CleanupParams struct {
	fx.In

	LC            fx.Lifecycle
	GlobalLimiter *RateLimiter `name:"globalLimiter"`
	LoginLimiter  *RateLimiter `name:"loginLimiter"`
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*Client),
		limit:   limit,
		window:  window,
	}
}

func NewGlobalRateLimiter() *RateLimiter {
	return NewRateLimiter(100, 1*time.Minute)
}

func NewLoginRateLimiter() *RateLimiter {
	return NewRateLimiter(20, 1*time.Minute)
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getClientKey(c)
		now := time.Now()

		rl.mu.Lock()
		client, exists := rl.clients[key]
		if !exists || now.After(client.ResetAt) {
			rl.clients[key] = &Client{
				Requests: 1,
				ResetAt:  now.Add(rl.window),
			}
			rl.mu.Unlock()
			c.Next()
			return
		}

		client.Requests++
		if client.Requests > rl.limit {
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		rl.mu.Unlock()
		c.Next()
	}
}

func getClientKey(c *gin.Context) string {
	authCookie, err := c.Cookie("access_token")

	if err != nil || authCookie == "" {
		return c.ClientIP()
	}

	token, _, err := new(jwt.Parser).ParseUnverified(authCookie, jwt.MapClaims{})
	if err != nil {
		return c.ClientIP()
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userId, exists := claims["user_id"]; exists {
			if idStr, ok := userId.(string); ok {
				return idStr
			}
		}
	}

	return c.ClientIP()
}

func StartRateLimitCleanup(p CleanupParams) {
	ticker := time.NewTicker(5 * time.Minute)
	stop := make(chan struct{})

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				for {
					select {
					case <-ticker.C:
						p.GlobalLimiter.Cleanup()
						p.LoginLimiter.Cleanup()
					case <-stop:
						ticker.Stop()
						return
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			close(stop)
			return nil
		},
	})
}

func (rl *RateLimiter) Cleanup() {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for key, client := range rl.clients {
		if now.After(client.ResetAt.Add(rl.window)) {
			delete(rl.clients, key)
		}
	}
}
