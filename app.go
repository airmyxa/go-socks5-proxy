package main

import (
	"fmt"

	"github.com/airmyxa/go-socks5-proxy/proto"
)

type Application struct {
	listener proto.ConnListener
}

func (a *Application) Run() {
	fmt.Println("Starting application")
	err := a.listener.Listen()
	if err != nil {
		fmt.Printf("Error while listening to connections: %s", err.Error())
	}
}
