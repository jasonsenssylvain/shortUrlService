package shorturllib

import (
	"net/http"
)

//该接口，最终实现的是短链接的processor和长连接的processor
type Processor interface {
	ProcessRequest(method, url string, params map[string]string, body []byte, w http.ResponseWriter, r *http.Request) error
}

type BaseProcessor struct {
	RedisCli      *RedisAdapter
	Configure     *Configure
	HostName      string
	LRU           *UrlLRU
	CountFunction CreateCountFunc
}
