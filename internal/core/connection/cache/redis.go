package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

type connection struct {
	client *redis.Client
	ctx    context.Context
}

// Configuration config redis connection
type Configuration struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

// Init init a new redis connection
func Init(cf *Configuration) error {
	addr := fmt.Sprintf("%s:%d", cf.Host, cf.Port)
	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: cf.Username,
		Password: cf.Password,
		DB:       cf.DB,
	})

	if err := conn.Ping(context.TODO()).Err(); err != nil {
		return err
	}

	client = conn

	return nil
}

// New new client connection
func New() *connection {
	return &connection{
		client: client,
		ctx:    context.Background(),
	}
}

// Service service interface
type Service interface {
	Set(key string, value interface{}, expiredTime time.Duration) error
	Get(key string, value interface{}) error
	GetKeys(pattern string) ([]string, error)
	Delete(key string) error
}

func (c *connection) Set(key string, value interface{}, expiredTime time.Duration) error {
	data, errMar := sonic.Marshal(&value)
	if errMar != nil {
		return errMar
	}

	err := c.client.Set(c.ctx, key, data, expiredTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *connection) Get(key string, value interface{}) error {
	val, err := c.client.Get(c.ctx, key).Result()
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

func (c *connection) GetKeys(pattern string) ([]string, error) {
	var keys []string
	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()
	for iter.Next(c.ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return keys, errors.New("key does not exists")
	}

	return keys, nil
}

func (c *connection) Delete(key string) error {
	err := c.client.Del(c.ctx, key).Err()
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

	err = c.client.Set(c.ctx, key, b.Bytes(), expiredTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) Get(key string, value interface{}) error {
	val, err := c.client.Get(c.ctx, key).Result()
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
