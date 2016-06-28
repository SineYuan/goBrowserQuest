package main

import (
	"fmt"
	"flag"
	"log"
	"github.com/kataras/iris"

	gs "github.com/SineYuan/goBrowserQuest/bqs"
)

var confFilePath = flag.String("config", "./config.json", "configuration file path")
var clientDir = flag.String("client", "", "BrowserQuest root to serve if provided")
var clientReqPrefix = flag.String("prefix", "/game", "request url prefix when client is provided, cannot be '/' ")

func main() {
	flag.Parse()
	config, err := gs.LoadConf(*confFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config)
	bqs := gs.NewBQS(config)

	if *clientDir != "" {
		iris.Static(*clientReqPrefix, *clientDir, 1)
	}
	iris.Any("/", bqs.ToIrisHandler())

	addr := fmt.Sprintf("%v:%v", config.Host, config.Port)
	log.Println("Server is running at " + addr)
	iris.Listen(addr)
}