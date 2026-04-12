package main

import (
	"log"

	"github.com/fullstop113/go-web3-demo/model"
	"github.com/fullstop113/go-web3-demo/router"
)

func main() {
	model.InitDB()

	r := router.InitRouter()

	log.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}
