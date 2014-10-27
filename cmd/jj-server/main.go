package main

import (
	"jj/server"
)

func main() {
	s := server.NewServer(":9999")
	s.Run()
}
