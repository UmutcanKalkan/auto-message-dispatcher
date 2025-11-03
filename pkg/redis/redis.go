package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewRedisClient(addr, password string, db int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{Client: client}, nil
}

func (c *Client) CacheSentMessage(ctx context.Context, messageID string, sentAt time.Time) error {
	key := fmt.Sprintf("sent_message:%s", messageID)
	value := sentAt.Format(time.RFC3339)

	err := c.Set(ctx, key, value, 7*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to cache sent message: %w", err)
	}

	return nil
}

func (c *Client) GetSentMessageTime(ctx context.Context, messageID string) (*time.Time, error) {
	key := fmt.Sprintf("sent_message:%s", messageID)

	val, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cached message: %w", err)
	}

	sentAt, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sent time: %w", err)
	}

	return &sentAt, nil
}
