package shorturllib

type UrlLRU struct {
	shortUrlLRU  *LRU
	originUrlLRU *LRU
	redisAdapter *RedisAdapter
}

func NewUrlLRU(redisAdapter *RedisAdapter) *UrlLRU {
	urlLru := &UrlLRU{}
	shortUrlLru, _ := NewLRU(1000, nil, urlLru.onGet)
	originUrlLru, _ := NewLRU(1000, nil, urlLru.onGet)

	urlLru.shortUrlLRU = shortUrlLru
	urlLru.originUrlLRU = originUrlLru
	urlLru.redisAdapter = redisAdapter
	return urlLru
}

func (t *UrlLRU) GetShortUrl(key interface{}) (string, bool) {
	v, ok := t.shortUrlLRU.Get(key)
	if !ok {
		return "", ok
	}
	return v.(string), true
}

func (t *UrlLRU) GetOriginUrl(key interface{}) (string, bool) {
	v, ok := t.originUrlLRU.Get(key)
	if !ok {
		return "", ok
	}
	return v.(string), true
}

func (t *UrlLRU) SetUrl(originalUrl, shortUrl string) bool {
	t.originUrlLRU.Add(shortUrl, originalUrl)
	t.shortUrlLRU.Add(originalUrl, shortUrl)
	return true
}

func (t *UrlLRU) onGet(key interface{}) (interface{}, error) {
	return t.redisAdapter.GetURL(key.(string))
}
