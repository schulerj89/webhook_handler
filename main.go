package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Payload struct {
	Request string `json:"request"`
}

var ctx = context.Background()

func main() {
	r := gin.Default()

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-pubsub:6379", // or the address of your Redis server
		Password: "",                  // no password set
		DB:       0,                   // use default DB
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/webhook", func(c *gin.Context) {
		var request Payload

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Publish the received data to a Redis channel
		err := redisClient.Publish(ctx, "webhookChannel", request.Request).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish to Redis"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User data received", "user": request})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
