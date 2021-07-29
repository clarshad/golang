package main

import (
	"github.com/clarshad/golang/terraform-as-service/server"
)

const port = 8080

func main() {
	server.Start(port)
}
