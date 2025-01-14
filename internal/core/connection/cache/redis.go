package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

// Client redis
type Connection struct {
	session *redis.Client
}

var (
	ctx     = context.Background()
	ctxTodo = context.TODO()
	client  = &Connection{}
)

// Configuration config Redis connection
type Configuration struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

func Init(cf *Configuration) error {
	dsn := fmt.Sprintf("%s:%d", cf.Host, cf.Port)

	conn := redis.NewClient(&redis.Options{
		Addr:         dsn,
		Username:     cf.Username,
		Password:     cf.Password,
		DB:           cf.DB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	if err := conn.Ping(ctxTodo).Err(); err != nil {
		return err
	}

	client = &Connection{
		session: conn,
	}

	return nil
}

// GetConnection get client connection
func GetConnection() *Connection {
	return client
}

// Service service interface
type Service interface {
	Set(key string, value interface{}, expiredTime time.Duration) error
	Get(key string, value interface{}) error
	GetKeys(pattern string) ([]string, error)
	Delete(key string) error
}

func (c *Connection) Set(key string, value interface{}, expiredTime time.Duration) error {
	data, errMar := sonic.Marshal(&value)
	if errMar != nil {
		return errMar
	}

	err := c.session.Set(ctx, key, data, expiredTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) Get(key string, value interface{}) error {
	val, err := c.session.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("key does not exists")
		}

		return err
	}

	errMar := sonic.Unmarshal([]byte(val), &value)
	if errMar != nil {
		return errMar
	}

	return nil
}

func (c *Connection) GetKeys(pattern string) ([]string, error) {
	var keys []string
	iter := c.session.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return keys, errors.New("key does not exists")
	}

	return keys, nil
}

func (c *Connection) Delete(key string) error {
	err := c.session.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

/*func (c *Connection) Set(key string, value interface{}, expiredTime time.Duration) error {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(value)
	if err != nil {
		return err
	}

	err = c.session.Set(ctx, key, b.Bytes(), expiredTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) Get(key string, value interface{}) error {
	val, err := c.session.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("key does not exists")
		}
		return err
	}

	b := bytes.Buffer{}
	b.Write([]byte(val))
	d := gob.NewDecoder(&b)
	err = d.Decode(value)
	if err != nil {
		return err
	}

	return nil
}*/
