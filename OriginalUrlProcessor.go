package main

import (
	"errors"
	"io"
	"net/http"

	"encoding/json"

	"github.com/jasoncodingnow/shortUrlService/shorturllib"
)

type OriginalUrlProcessor struct {
	*shorturllib.BaseProcessor
	CountChannel chan shorturllib.CountChannel
}

const POST string = "POST"
const TOKEN string = "token"
const ORIGINAL_URL string = "original"
const SHORT_URL string = "short"

func (t *OriginalUrlProcessor) ProcessRequest(method, url string, params map[string]string, body []byte, w http.ResponseWriter, r *http.Request) error {
	if method != POST {
		return errors.New("method must be POST")
	}
	var bodyInfo map[string]interface{}
	err := json.Unmarshal(body, &bodyInfo)
	if err != nil {
		return err
	}

	_, tokenExists := bodyInfo[TOKEN].(string)
	originalUrl, originalUrlExists := bodyInfo[ORIGINAL_URL].(string)

	if !tokenExists || !originalUrlExists {
		return errors.New("Post error")
	}
	shortUrl, err := t.createUrl(originalUrl)
	if err != nil {
		return err
	}

	response, err := t.createResponseJson(shortUrl)
	if err != nil {
		return err
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")
	header.Add("charset", "UTF-8")
	io.WriteString(w, response)

	return nil
}

func (t *OriginalUrlProcessor) createUrl(originalUrl string) (string, error) {
	short, ok := t.LRU.GetShortUrl(originalUrl)
	if ok {
		return short, nil
	}

	count, err := t.CountFunction()
	if err != nil {
		return "", err
	}

	shortUrl, err := shorturllib.TransNumToString(count)
	if err != nil {
		return "", err
	}

	t.LRU.SetUrl(originalUrl, shortUrl)
	return shortUrl, nil
}

func (t *OriginalUrlProcessor) createResponseJson(shortUrl string) (string, error) {
	jsonRes := make(map[string]interface{})
	jsonRes[shortUrl] = t.HostName + shortUrl

	res, err := json.Marshal(jsonRes)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
