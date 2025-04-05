package limiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBucketLimiter struct {
	redisClient *redis.Client
	maxTokens   int     // Maximum tokens the bucket can hold
	refillRate  float64 // Tokens per second
	window      time.Duration
}

func NewTokenBucketLimiter(redisClient *redis.Client, maxTokens int, refillRate float64, window time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		redisClient: redisClient,
		maxTokens:   maxTokens,
		refillRate:  refillRate,
		window:      window,
	}
}

func (t *TokenBucketLimiter) AllowRequest(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("token_bucket:%s", ip)
	now := time.Now().Unix()

	// Fetch current token count and last update time
	data, err := t.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return false, err
	}

	var tokens int
	var lastUpdated int64
	if len(data) == 0 {
		// First request: set max tokens
		tokens = t.maxTokens
		lastUpdated = now
	} else {
		// Read current tokens and last update timestamp
		tokens = stringToInt(t.redisClient.HGet(ctx, key, "tokens").Val())
		lastUpdated = int64(stringToInt(t.redisClient.HGet(ctx, key, "last_updated").Val()))
	}

	// Calculate new tokens based on elapsed time
	elapsed := now - lastUpdated
	newTokens := int(float64(elapsed) * t.refillRate)
	if newTokens > 0 {
		tokens += newTokens
		if tokens > t.maxTokens {
			tokens = t.maxTokens
		}
		lastUpdated = now
	}

	// If no tokens left, deny request
	if tokens <= 0 {
		return false, nil
	}

	// Consume 1 token
	tokens--

	// Store updated values in Redis
	t.redisClient.HSet(ctx, key, map[string]any{
		"tokens":       tokens,
		"last_updated": lastUpdated,
	})
	t.redisClient.Expire(ctx, key, t.window) // Set TTL to auto-remove old records

	return true, nil
}

func stringToInt(snum string) int {
	num, _ := strconv.Atoi(snum)
	return num
}
