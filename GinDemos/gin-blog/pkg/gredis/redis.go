package gredis

import (
	"encoding/json"
	"time"

	"gin-blog/pkg/setting"

	"github.com/gomodule/redigo/redis"
)

/*
备注说明：
1、设置 RedisConn 为 redis.Pool（连接池）并配置了它的一些参数：
	Dial：提供创建和配置应用程序连接的一个函数
	TestOnBorrow：可选的应用程序检查健康功能
	MaxIdle：最大空闲连接数
	MaxActive：在给定时间内，允许分配的最大连接数（当为零时，没有限制）
	IdleTimeout：在给定时间内将会保持空闲状态，若到达时间限制则关闭连接（当为零时，没有限制）
2、封装基础方法

	文件内包含 Set、Exists、Get、Delete、LikeDeletes 用于支撑目前的业务逻辑，而在里面涉及到了如方法：
	（1）RedisConn.Get()：在连接池中获取一个活跃连接
	（2）conn.Do(commandName string, args ...interface{})：向 Redis 服务器发送命令并返回收到的答复
	（3）redis.Bool(reply interface{}, err error)：将命令返回转为布尔值
	（4）redis.Bytes(reply interface{}, err error)：将命令返回转为 Bytes
	（5）redis.Strings(reply interface{}, err error)：将命令返回转为 []string
	在 redigo 中包含大量类似的方法，万变不离其宗，建议熟悉其使用规则和 Redis命令 即可
*/

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
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

	return nil
}

func Set(key string, data interface{}, time int) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	reply, err := redis.Bool(conn.Do("SET", key, value))
	conn.Do("EXPIRE", key, time)

	return reply, err
}

func Exist(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
