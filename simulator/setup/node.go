package functions

import (
	"os"

	"github.com/ermes-labs/api-go/api"
	"github.com/ermes-labs/api-go/infrastructure"
	rc "github.com/ermes-labs/storage-redis/packages/go"
	"github.com/redis/go-redis/v9"
)

// The node that the function is running on.
var Node *api.Node

// The Redis client.
var redisClient *redis.Client

func init() {
	// Get the node from the environment variable.
	jsonNode := os.Getenv("ERMES_NODE")
	// Unmarshal the environment variable to get the node.
	infraNode, err := infrastructure.UnmarshalNode([]byte(jsonNode))
	// Check if there was an error unmarshalling the node.
	if err != nil {
		panic(err)
	}

	// Get the Redis connection details from the environment variables.
	redisHost := envOrDefault("REDIS_HOST", "localhost")
	redisPort := envOrDefault("REDIS_PORT", "6379")
	redisPassword := envOrDefault("REDIS_PASSWORD", "")
	// Create a new Redis client.
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	// The Redis commands.
	var RedisCommands = rc.NewRedisCommands(redisClient)
	// Create a new node with the Redis commands.
	Node = api.NewNode(*infraNode, RedisCommands)
}

// Get the value of an environment variable or return a default value.
func envOrDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}
