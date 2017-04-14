package shorturllib

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type RedisAdapter struct {
	conn   redis.Conn
	config *Configure
}

const SHORT_URL_COUNT_KEY string = "short_url_count"

//NewRedisAdapter 实例化redis的操作实例
func NewRedisAdapter(config *Configure) (*RedisAdapter, error) {
	redisCli := &RedisAdapter{}
	redisCli.config = config

	redisHost, _ := config.GetRedisHost()
	redisPort, _ := config.GetRedisPort()

	connStr := fmt.Sprintf("%v:%v", redisHost, redisPort)
	conn, err := redis.Dial("tcp", connStr)
	if err != nil {
		return nil, err
	}

	redisCli.conn = conn
	return redisCli, nil
}

//Release 释放连接
func (t *RedisAdapter) Release() {
	t.conn.Close()
}

//InitCountService 初始化短链接计数服务
func (t *RedisAdapter) InitCountService() error {
	_, err := t.conn.Do("INCR", SHORT_URL_COUNT_KEY)
	if err != nil {
		return err
	}
	return nil
}

//NewShortURLCount 增加短链接计数
func (t *RedisAdapter) NewShortURLCount() (int64, error) {
	count, err := redis.Int64(t.conn.Do("INCR", SHORT_URL_COUNT_KEY))
	if err != nil {
		return 0, err
	}
	return count, nil
}

//SetURL 保存短链接
func (t *RedisAdapter) SetURL(originalURL, shortURL string) error {
	key := fmt.Sprintf("short:%v", shortURL)
	_, err := t.conn.Do("SET", key, originalURL)
	if err != nil {
		return err
	}
	return nil
}

//GetURL 根据短连接返回原始链接
func (t *RedisAdapter) GetURL(shortURL string) (string, error) {
	key := fmt.Sprintf("short:%v", shortURL)
	originalURL, err := redis.String(t.conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return originalURL, nil
}
