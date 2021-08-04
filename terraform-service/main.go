package main

import (
	"github.com/clarshad/golang/terraform-service/server"
)

const port = 8080

func main() {
	server.Handle(port)
}
