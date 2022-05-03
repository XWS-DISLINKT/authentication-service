package main

import (
	"authentication_service/startup"
	"authentication_service/startup/config"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
