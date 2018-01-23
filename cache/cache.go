package cache

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)
const (
	server                 = ":6379"
)

var Cache = GetCache()

func GetCache() *RedisStorage {

	return &RedisStorage{
		Pool:      NewPool(),
	}

}

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				panic(err.Error())
			}
			log.Printf("Dail master redis server %s succeefully!", server)
			return c, err
		},
	}
}

type RedisStorage struct {
	Pool *redis.Pool
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
}

func (r *RedisStorage) Set(key string, value string, expiration time.Duration) error {

	var err error
	c:= r.Pool.Get()
	_, err = c.Do("SET", key, value)
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

