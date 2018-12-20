package tools

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
)

// redis 队列相关函数

var (
	pool *redis.Pool
)

// RedisConf redis 配子项
type RedisConf struct {
	Server      string `yaml:"server"`
	Passwd      string `yaml:"passwd"`
	DB          int64  `yaml:"db"`
	MaxIdle     int    `yaml:"maxIdle"`
	MaxActive   int    `yaml:"maxActive"`
	IdleTimeout int64  `yaml:"idleTimeOut"`
}

//RedisPoolInit redis初始化
func RedisPoolInit(conf *RedisConf) {
	server := conf.Server
	passwd := conf.Passwd
	db := conf.DB
	maxActive := conf.MaxActive
	maxIdle := conf.MaxIdle
	idleTimeout := conf.IdleTimeout

	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if passwd != "" {
				if _, err := c.Do("AUTH", passwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			if db != 0 {
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// Enqueue 入队
func Enqueue(b []byte, queue string) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("RPUSH", queue, b)
	return err
}

//Dequeue 出队操作
func Dequeue(queue string) ([]byte, error) {

	count, err := queueCount(queue)
	if err != nil || count == 0 {
		return nil, errors.New("no job to do")
	}
	glog.V(2).Info("[queueCount()] Jobs count:", count)

	c := pool.Get()
	defer c.Close()

	r, err := redis.Bytes(c.Do("LPOP", queue))
	return r, err

}

//queueCount 队列数量
func queueCount(queue string) (int, error) {
	c := pool.Get()
	defer c.Close()

	lenqueue, err := c.Do("LLEN", queue)
	if err != nil {
		return 0, err
	}

	count, ok := lenqueue.(int64)
	if !ok {
		return 0, errors.New("获取数量类型转换错误!")
	}
	return int(count), nil
}

//GetStringValue 获取一个缓存值,如果不存在,则返回err
func GetStringValue(key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))

}

//GetByteValue 获取一个[]byte类型的值
func GetByteValue(key string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()
	r, err := redis.Bytes(conn.Do("GET", key))
	return r, err
}

//SetValue 设置缓存值
func SetValue(key string, value string, exp int) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	if exp > 0 {
		_, err = conn.Do("EXPIRE", key, exp)
		if err != nil {
			return err
		}
	}
	return nil
}

//DelKey 删除一个缓存
func DelKey(key string) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}
