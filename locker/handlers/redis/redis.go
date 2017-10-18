package redis

import (
	"errors"
	"time"

	"github.com/digota/digota/storage/object"

	"github.com/digota/digota/config"
	"github.com/garyburd/redigo/redis"
)

type locker struct {
	rp *redis.Pool
}

const separator = "."

// NewLocker return new redis based lock
func NewLocker(lockerConfig config.Locker) (*locker, error) {
	if len(lockerConfig.Address) < 1 {
		return nil, errors.New("No redis address provided")
	}
	p := newPool(lockerConfig.Address[0], "")
	return &locker{p}, nil
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		MaxActive:   20,
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

func getKey(doc object.Interface) (string, error) {
	if doc.GetNamespace() == "" || doc.GetId() == "" {
		return "", errors.New("Obj is missing information to make that lock")
	}
	return doc.GetNamespace() + separator + doc.GetId(), nil
}

func (l *locker) Close() error {
	return l.rp.Close()
}

func (l *locker) Lock(doc object.Interface) (func() error, error) {
	key, err := getKey(doc)
	if err != nil {
		return nil, err
	}

	conn := l.rp.Get()
	_, err = redis.String(conn.Do("SET", key, "NX"))
	conn.Close()
	if err != nil {
		return nil, err
	}

	return func() error { return l.unlock(key) }, nil
}

func (l *locker) TryLock(doc object.Interface, t time.Duration) (func() error, error) {
	key, err := getKey(doc)
	if err != nil {
		return nil, err
	}

	ch := make(chan error)
	conn := l.rp.Get()
	defer conn.Close()

	go func(c redis.Conn) {
		_, err = redis.String(c.Do("SET", key, "NX"))
		select {
		case ch <- err:
		default:
		}
	}(conn)

	select {
	case err = <-ch:
		if err != nil {
			return nil, err
		}
		return func() error { return l.unlock(key) }, nil
	case <-time.After(t):
		return nil, errors.New("Timeout reached")
	}
}

func (l *locker) unlock(key string) error {
	conn := l.rp.Get()
	_, err := conn.Do("DEL", key)
	conn.Close()
	return err
}
