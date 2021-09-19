package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"time"
)

var redisClient *redis.Client

func main() {

	ctx := context.Background()
	redisClient.Set(ctx, "a", RandStringBytesMaskImpr(100), time.Second*100)
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "",
		PoolSize: 3,
	})
	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err == redis.Nil {
		log.Fatal("Redis异常", err)
	} else if err != nil {
		log.Fatal("失败:", err.Error())
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
