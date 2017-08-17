package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gerifield/digota/config"
)

type lock struct {
	rp *redis.Pool
}

// NewLocker return new redis based lock
func NewLocker(lockerConfig config.Locker) (*lock, error) {
	p := newPool(lockerConfig.Address, "")
	return &lock{c}, nil
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     20,
		MaxActive:   80, // 10k is the max currently on the redis server
		IdleTimeout: 240 * time.Second,
		Wait:        true, // Wait for the connection pool, no connection pool exhausted error
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialConnectTimeout(100*time.Millisecond), redis.DialReadTimeout(200*time.Millisecond), redis.DialWriteTimeout(200*time.Millisecond))
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
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
