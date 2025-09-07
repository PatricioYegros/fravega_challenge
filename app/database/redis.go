package database

import (
	"challenge_pyegros/app/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Could not connect to Redis:", err)
		return nil
	}
	fmt.Println("Connected to Redis:", pong)
	return rdb
}

func SetEventDataFromRedis(key string, eventResponse *models.ResponseUpdate, rdb *redis.Client) error {
	value, err := json.Marshal(eventResponse)
	if err != nil {
		return err
	}

	err = rdb.Set(context.Background(), "event: "+key, value, 24*time.Hour).Err()
	return err
}

func GetEventDataFromRedis(key string, rdb *redis.Client) (*models.ResponseUpdate, error) {
	val, err := rdb.Get(context.Background(), "event: "+key).Bytes()
	if err != nil {
		return nil, err
	}

	var response models.ResponseUpdate
	if err = json.Unmarshal(val, &response); err != nil {
		return nil, err
	}
	return &response, nil

}

func GetOrderDataFromRedis(key string, rdb *redis.Client) (*models.ResponseCreate, error) {
	val, err := rdb.Get(context.Background(), "order: "+key).Bytes()
	if err != nil {
		return nil, err
	}

	var response models.ResponseCreate
	if err = json.Unmarshal(val, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func SetOrderDataFromRedis(key string, orderResponse *models.ResponseCreate, rdb *redis.Client) error {
	value, err := json.Marshal(orderResponse)
	if err != nil {
		return err
	}

	err = rdb.Set(context.Background(), "order: "+key, value, 24*time.Hour).Err()
	return err
}
