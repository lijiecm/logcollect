package main

import (
	"github.com/weihualiu/logcollect/server"
	"log"
	"os"
)

func main() {
	//logset()
	log.Printf("logcollect is starting......")

	server.Start()
}

func logset() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file:%v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetFlags(log.Flags() | log.Lshortfile)
}
