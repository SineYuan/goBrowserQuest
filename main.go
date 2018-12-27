package main

import (
	"fmt"
	"flag"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	gs "github.com/SineYuan/goBrowserQuest/bqs"
)

var confFilePath = flag.String("config", "./config.json", "configuration file path")
var clientDir = flag.String("client", "", "BrowserQuest root directory to serve if provided")
var clientReqPrefix = flag.String("prefix", "/game", "request url prefix when client is provided, cannot be '/' ")

func main() {
	flag.Parse()
	config, err := gs.LoadConf(*confFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config)
	bqs := gs.NewBQS(config)

	e := echo.New()
	e.Use(middleware.Recover())

	if *clientDir != "" {
		e.Static(*clientReqPrefix, *clientDir)
	}
	e.Any("/", bqs.ToEchoHandler())

	addr := fmt.Sprintf("%v:%v", config.Host, config.Port)
	log.Println("Server is running at " + addr)
	e.Logger.Fatal(e.Start(addr))
}