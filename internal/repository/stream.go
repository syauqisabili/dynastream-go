package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"stream-session-api/domain"
	"stream-session-api/internal/conf/network"
	"time"

	"github.com/redis/go-redis/v9"
)

type streamRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewStream() domain.StreamRepository {
	addr := fmt.Sprintf("%s:%d", network.Get().Redis.Ip, network.Get().Redis.Port)
	password := network.Get().Redis.Password
	db := network.Get().Redis.DatabaseIndex

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password, // no password set
		DB:           int(db),
		DialTimeout:  5 * time.Second, // Wait to conenct
		ReadTimeout:  5 * time.Second, // Wait to read
		WriteTimeout: 5 * time.Second, // Wait to get
	})

	return &streamRepository{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *streamRepository) Close() {
	r.client.Close()
}

func (r *streamRepository) GetAll() ([]domain.Stream, error) {
	var cursor uint64
	var results []domain.Stream

	for {
		// Scan for matching keys
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(r.ctx, cursor, "log:stream:*", 0).Result()
		if err != nil {
			return nil, err
		}

		// Fetch values for the keys
		for _, key := range keys {
			// Get the value for each key
			value, err := r.client.Get(r.ctx, key).Result()
			if err != nil {
				return nil, err
			}

			// Unmarshal the JSON into the struct
			result := &domain.Stream{}

			if err := json.Unmarshal([]byte(value), result); err != nil {
				return nil, err
			}
			// Add the result to the results slice
			results = append(results, *result)
		}

		// Break if cursor is 0 (no more keys)
		if cursor == 0 {
			break
		}
	}

	return results, nil
}

func (r *streamRepository) FindByUuid(uuid string) *domain.Stream {
	// Get the value for each key
	value, err := r.client.Get(r.ctx, fmt.Sprintf("log:stream:%s", uuid)).Result()
	if err != nil {
		return nil
	}

	// Unmarshal the JSON into the struct
	var result *domain.Stream
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return nil
	}
	if result.Uuid == uuid {
		return result
	}

	return nil
}

func (r *streamRepository) Insert(stream *domain.Stream) error {
	// Marshal the struct into JSON
	json, err := json.Marshal(stream)
	if err != nil {
		return err
	}

	// Save data JSON to redis
	err = r.client.Set(r.ctx, fmt.Sprintf("log:stream:%s", stream.Uuid), json, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *streamRepository) Delete(uuid string) error {
	key := fmt.Sprintf("log:stream:%s", uuid)
	if err := r.client.Del(r.ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
