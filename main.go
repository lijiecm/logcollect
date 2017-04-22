package main

import (
	"github.com/weihualiu/logcollect/server"
	"log"
)

func main() {
	log.Printf("logcollect is starting......")
	server.Start()
}
