package main

import "github.com/luizrgf2/universal-flow/internal/presentation"

func main() {
	route := presentation.StartServer()
	route.Run(":8080")
}
