package main

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

func TestGetUser(t *testing.T) {

	rclnt := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	val, err := rclnt.Get("bjamesdowning@gmail.com").Result()
	if err != nil {
		fmt.Println("Get:", err)
	}
	fmt.Println("KEY:", val)

}
