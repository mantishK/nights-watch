package config

import (
	"time"

	"plivo/nights-watch/log"

	"github.com/garyburd/redigo/redis"
)

var RedisPool redis.Pool

func init() {
	server, ok := GetString("redis_ip")
	if !ok {
		log.Err("Config values missing")
		panic("Config values for redis missing, unable to start the server")
	}
	pass, ok := GetString("redis_pass")
	if !ok {
		log.Err("Config values missing")
		panic("Config values for redis missing, unable to start the server")
	}
	RedisPool = redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(pass) != 0 {
				if _, err := c.Do("AUTH", pass); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
