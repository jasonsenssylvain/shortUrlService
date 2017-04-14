package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	lib "github.com/jasoncodingnow/shortUrlService/shorturllib"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "conf", "config.ini", "configure file full path")
	flag.Parse()

	fmt.Printf("[INFO] Read config file...\n")
	fmt.Printf("config file path is %s", configFile)

	config, err := lib.NewConfigure(configFile)
	if err != nil {
		fmt.Printf("[ERROR] Parse configure file error: %v\n", err)
	}

	fmt.Printf("[INFO] Start Redis Client... \n")
	redisCli, err := lib.NewRedisAdapter(config)
	if err != nil {
		fmt.Printf("[ERROR] Redis init failed...")
		return
	}

	if config.GetRedisStatus() {
		err = redisCli.InitCountService()
		if err != nil {
			fmt.Printf("[ERROR] Init Redis key count failed... " + err.Error())
		}
	}

	countChannel := make(chan lib.CountChannel, 1000)
	go CountThread(countChannel)

	counterFunc := lib.CreateCounter(countChannel, redisCli)

	urlLRU := lib.NewUrlLRU(redisCli)
	if err != nil {
		fmt.Printf("[ERROR]LRU init fail...\n")
	}

	baseProcessor := &lib.BaseProcessor{redisCli, config, config.GetHostInfo(), urlLRU, counterFunc}
	original := &OriginalUrlProcessor{baseProcessor, countChannel}
	short := &ShortUrlProcessor{baseProcessor}

	router := &lib.Router{map[int]lib.Processor{
		0: short,
		1: original,
	}}

	port, _ := config.GetPort()
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("[INFO]Service Starting addr :%v,port :%v\n", addr, port)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		//logger.Error("Server start fail: %v", err)
		os.Exit(1)
	}
}

func CountThread(countChan chan lib.CountChannel) {
	var count int64
	count = 1000
	for {
		select {
		case ok := <-countChan:
			count = count + 1
			fmt.Println("count add, now is " + strconv.FormatInt(count, 10))
			ok.CountOutChan <- count
		}
	}
}
