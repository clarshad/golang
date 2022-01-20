package main

import (
	"github.com/clarshad/golang/terraform-service/server"
)

// defualt port to listen
const port = 8080

func main() {
	server.Handle(port)
}
