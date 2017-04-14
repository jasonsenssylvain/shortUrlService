package main

import (
	"errors"
	"net/http"

	"strings"

	"github.com/jasoncodingnow/shortUrlService/shorturllib"
)

type ShortUrlProcessor struct {
	*shorturllib.BaseProcessor
}

func (t *ShortUrlProcessor) ProcessRequest(method, url string, params map[string]string, body []byte, w http.ResponseWriter, r *http.Request) error {
	hostname := t.HostName
	urls := strings.Split(url, hostname+"/")
	var key string
	for _, url := range urls {
		key = url[1:]
	}
	originalUrl, err := t.GetOriginalUrl(key)
	if err != nil {
		return err
	}
	http.Redirect(w, r, originalUrl, http.StatusMovedPermanently)
	return nil
}

func (t *ShortUrlProcessor) GetOriginalUrl(requestUrl string) (string, error) {
	originalUrl, ok := t.LRU.GetOriginUrl(requestUrl)
	if !ok {
		return "", errors.New("cannot find url")
	}

	return originalUrl, nil
}
