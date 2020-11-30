package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
)

// Config contains the configuration for the app.
type Config struct {
	Port          int    `envconfig:"PORT"`
	RedisHost     string `envconfig:"REDIS_HOST"`
	RedisPort     int    `envconfig:"REDIS_PORT"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
	RedisDb       int    `envconfig:"REDIS_DB"`
}
var ctx = context.Background()
// Count is the response object.
type Count struct {
	Value int `json:"count"`
}

func main() {
	var c Config
	err := envconfig.Process("app", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),//fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort),
		Password: c.RedisPassword,
		DB:       c.RedisDb,
	})

	//pong, err := r.Ping(ctx).Result()
	//if err == nil {
	//	log.Fatal(pong, err)
	//} else {
	//	log.Fatal(err)
	//}


	log.Println(c.RedisHost, c.RedisPort, c.RedisPassword, c.RedisDb)
	router := gin.Default()

	router.GET("/counter/:id", func(c *gin.Context) {
		ctx := context.Background()
		id := c.Param("id")
		val, err := r.Get(ctx, id).Result()
		if err != nil {
			if err == redis.Nil {
				c.JSON(http.StatusOK, Count{Value: 1})
				if err := r.Set(ctx, id, 1, 0).Err(); err != nil {
					log.Println("failed to update cache")
				}
				return
			}
			c.String(http.StatusInternalServerError, "failed to get count")
		}
		ct, err := strconv.Atoi(val)
		if err != nil {
			log.Println("failed to parse value")
		}

		c.JSON(http.StatusOK, Count{Value: ct + 1})
		if err := r.Set(ctx, id, ct+1, 0).Err(); err != nil {
			log.Println("failed to update cache")
		}
		return
	})

	router.GET("/ping", func(c *gin.Context) {
		if _, err := r.Ping(context.Background()).Result(); err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(http.StatusOK)
		}
	})

	_ = router.Run(fmt.Sprintf(":%d", c.Port))
}
