package cache

import (
	"github.com/gomodule/redigo/redis"
	"log"
)
const (
	server                 = ":6379"
)

var Cache *RedisStorage
var Pool *redis.Pool

func init() {
	Pool = NewPool()
	Cache = &RedisStorage{
		Pool: Pool,
	}
}

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				log.Println(err.Error())
			}
			return c, err
		},
	}
}

type RedisStorage struct {
	Pool *redis.Pool
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
}

func (r *RedisStorage) Set(key string, value string, expiration int32) error {

	var err error
	c:= r.Pool.Get()
	_, err = c.Do("SETEX", key, expiration, value)
	if err != nil {
		log.Printf("Could not SET %s:%s.", key, err)
		return err
	}

	return nil
}

func (r *RedisStorage) Get(key string) (string, error) {

	c:= r.Pool.Get()
	resp, err := redis.String(c.Do("GET", key))

	if err != nil {
		log.Printf("cache: Marshal key=%q failed: %s", key, err)
		return "", err
	}

	return resp, nil
}

func (r *RedisStorage) Del(key string) error {
	c:= r.Pool.Get()
	_, err := redis.String(c.Do("DEL", key))
	return err
}

