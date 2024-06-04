package main

// Import packages
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

// define context for use in redis
var ctx = context.Background()

// Generate Random string with a custom length
func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// Function to insert a map of key values into the redis
func RedisSetKey(redis_connection *redis.Client, keyValues map[string]string) {

	fmt.Println("Start insert keys")
	err := redis_connection.MSet(ctx, keyValues).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("All keys insert successfully")
}

// Function for Scan keys with a specific key name
func RedisScanKeys(redis_connection *redis.Client, keyName string) {
	startTime := time.Now()
	var cursor uint64
	var keys []string
	var err error
	pattern := fmt.Sprintf("*" + keyName + "*")
	for {
		keys, cursor, err = redis_connection.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			fmt.Println("Error while scanning:", err)
			return
		}

		for _, key := range keys {
			fmt.Println("Found key:", key)
		}

		if cursor == 0 {
			break
		}
	}
	fmt.Println(time.Since(startTime).Milliseconds())

}

// Function for Get specific key from redis
func RedisGetKeys(redis_connection *redis.Client, keyName string) {
	startTime := time.Now()

	key, err := redis_connection.Get(ctx, keyName).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Found key:", key)
	fmt.Println(time.Since(startTime).Milliseconds())

}

func main() {
	// Load .env file and get variables
	if err := godotenv.Load(); err != nil {
		panic("Failed to load environment variables")
	}
	key_count, _ := strconv.Atoi(os.Getenv("KEY_COUNT"))
	key_char_size, _ := strconv.Atoi(os.Getenv("KEY_CHAR_SIZE"))
	value_char_size, _ := strconv.Atoi(os.Getenv("VALUE_CHAR_SIZE"))
	run_count, _ := strconv.Atoi(os.Getenv("RUN_COUNT"))
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")

	// Create redis connection address
	redis_connection := fmt.Sprintf(redis_host + ":" + redis_port)

	// loop for run application for several times
	for i := 1; i <= run_count; i++ {
		fmt.Printf("************** Run %d ************\n", i)
		// Create connection for start insert keys
		rdb := redis.NewClient(&redis.Options{
			Addr:     redis_connection,
			Password: "",
			DB:       0,
		})
		fmt.Println("Connected to redis")
		// FlushAll keys from redis
		rdb.FlushAll(ctx)
		var testKey string

		// Generate a map of random key values for insert into redis
		randomMap := make(map[string]string, key_count)

		for i := 0; i < key_count; i++ {
			key := GenerateRandomString(key_char_size)
			value := GenerateRandomString(value_char_size)
			randomMap[key] = value
			testKey = key
		}

		// Call Redis Set Function and Close redis client connection
		RedisSetKey(rdb, randomMap)
		rdb.Close()

		// Create another connection for scan and get keys
		rdb2 := redis.NewClient(&redis.Options{
			Addr:     redis_connection,
			Password: "",
			DB:       0,
		})

		// Call Redis Scan Function to find a key
		fmt.Println("*************** Scan Method **************")
		RedisScanKeys(rdb2, testKey)

		// Call Redis Get Function to find a key
		fmt.Println("*************** Get Method ***************")
		RedisGetKeys(rdb2, testKey)

		// Close Connection
		rdb2.Close()

	}

}
