package main

import (
	"github.com/sounishnath003/url-shortner-service-golang/internal/core"
	"github.com/sounishnath003/url-shortner-service-golang/internal/server"
)

func main() {
	co := core.InitCore()

	server := server.NewServer(co)
	panic(server.Run())
}
