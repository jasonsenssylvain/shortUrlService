package shorturllib

import (
	"container/list"
	"fmt"
)

type CountChannel struct {
	Ok           int64
	CountOutChan chan int64
}

type CreateCountFunc func() (int64, error)

func CreateCounter(countChannel chan CountChannel, redisCli *RedisAdapter) CreateCountFunc {
	return func() (int64, error) {
		fmt.Println("try to add count")
		count, err := redisCli.NewShortURLCount()
		if err != nil {
			return 0, err
		}
		return count, nil
	}
}

func TransNumToString(num int64) (string, error) {
	var base int64
	base = 62
	baseHex := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	output_list := list.New()
	for num/base != 0 {
		output_list.PushFront(num % base)
		num = num / base
	}
	output_list.PushFront(num % base)
	str := ""
	for iter := output_list.Front(); iter != nil; iter = iter.Next() {
		str = str + string(baseHex[int(iter.Value.(int64))])
	}
	return str, nil
}
