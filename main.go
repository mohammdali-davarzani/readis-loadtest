package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func RedisSetKey(redis_connection *redis.Client, keyValues map[string]string) {

	// key := os.Getenv("KEYNAME")
	fmt.Println("Start insert keys")
	err := redis_connection.MSet(ctx, keyValues).Err()
	if err != nil {
		panic(err)
	}
}

func RedisScanKeys(redis_connection *redis.Client, keyName string) {
	startTime := time.Now()
	var cursor uint64
	var err error
	pattern := fmt.Sprintf("*" + keyName + "*")
	for {
		_, cursor, err = redis_connection.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			fmt.Println("Error while scanning:", err)
			return
		}

		if cursor == 0 {
			break
		}
	}
	fmt.Println(time.Since(startTime).Milliseconds())

}

func RedisGetKeys(redis_connection *redis.Client, keyName string) {
	startTime := time.Now()

	_, err := redis_connection.Get(ctx, keyName).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Since(startTime).Milliseconds())

}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load environment variables")
	}
	key_count, _ := strconv.Atoi(os.Getenv("KEY_COUNT"))
	key_char_size, _ := strconv.Atoi(os.Getenv("KEY_CHAR_SIZE"))
	value_char_size, _ := strconv.Atoi(os.Getenv("VALUE_CHAR_SIZE"))
	run_count, _ := strconv.Atoi(os.Getenv("RUN_COUNT"))
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_connection := fmt.Sprintf(redis_host + ":" + redis_port)
	for i := 1; i <= run_count; i++ {
		fmt.Printf("************** Run %d ************\n", i)
		rdb := redis.NewClient(&redis.Options{
			Addr:     redis_connection,
			Password: "",
			DB:       0,
		})
		fmt.Println("Connected to redis")
		rdb.FlushAll(ctx)
		var testKey string
		randomMap := make(map[string]string, key_count)

		for i := 0; i < key_count; i++ {
			key := GenerateRandomString(key_char_size)
			value := GenerateRandomString(value_char_size)
			randomMap[key] = value
			testKey = key
		}

		RedisSetKey(rdb, randomMap)
		rdb.Close()

		rdb2 := redis.NewClient(&redis.Options{
			Addr:     redis_connection,
			Password: "",
			DB:       0,
		})
		fmt.Println("*************** Scan Method **************")
		RedisScanKeys(rdb2, testKey)
		fmt.Println("*************** Get Method ***************")
		RedisGetKeys(rdb2, testKey)
		rdb2.Close()

	}

}
