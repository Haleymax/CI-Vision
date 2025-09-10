package main

import (
	"civ/cmd/server"
)

// main is the program entry point. It delegates initialization and execution of the application server to server.RunServer().
func main() {
	server.RunServer()
}
