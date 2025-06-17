package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

var (
	connection = &redis.Client{}
	ctx        = context.Background()
)

// Configuration config redis connection
type Configuration struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

type client struct {
	client *redis.Client
}

type Client interface {
	Set(key string, value interface{}, expiredTime time.Duration) error
	Get(key string, value interface{}) error
	GetKeys(pattern string) ([]string, error)
	Delete(key string) error
	Close() error
}

// Init init a new redis connection
func Init(cf *Configuration) error {
	addr := fmt.Sprintf("%s:%d", cf.Host, cf.Port)
	connection = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: cf.Username,
		Password: cf.Password,
		DB:       cf.DB,
	})

	err := connection.Ping(context.TODO()).Err()
	if err != nil {
		return err
	}

	return nil
}

// New new client connection
func New() Client {
	return &client{
		client: connection,
	}
}

func (c *client) Set(key string, value interface{}, expiredTime time.Duration) error {
	data, errMar := json.Marshal(&value)
	if errMar != nil {
		return errMar
	}

	err := c.client.Set(ctx, key, data, expiredTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Get(key string, value interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("key does not exists")
		}

		return err
	}

	errMar := json.Unmarshal([]byte(val), &value)
	if errMar != nil {
		return errMar
	}

	return nil
}

func (c *client) GetKeys(pattern string) ([]string, error) {
	var keys []string
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return keys, errors.New("key does not exists")
	}

	return keys, nil
}

func (c *client) Delete(key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// Close close connection
func (c *client) Close() error {
	err := c.client.Close()
	if err != nil {
		return err
	}

	return nil
}

/*func (c *client) Set(key string, value interface{}, expiredTime time.Duration) error {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(value)
	if err != nil {
		return err
	}

	err = c.client.Set(ctx, key, b.Bytes(), expiredTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Get(key string, value interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
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
