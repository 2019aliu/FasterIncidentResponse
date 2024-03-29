package redis

import (
	"github.com/gomodule/redigo/redis"
)

func InitCache() redis.Conn {
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	// Assign the connection to the package level `cache` variable
	cache := conn

	return cache
}
