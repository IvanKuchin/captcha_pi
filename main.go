package main

import (
	"log"

	server "github.com/ivankuchin/captcha_pi/server"
	configreader "github.com/ivankuchin/timecard.ru-api/config-reader"
)

func SetLogFlags() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	SetLogFlags()

	config, err := configreader.Read()
	if err != nil {
		log.Panic(err.Error())
	}

	log.Print(config.Listenport)

	server.SetConfig(*config)
	server.Run()
}
