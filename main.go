package main

import (
	"log"

	"github.com/fullstop113/go-web3-demo/config"
	"github.com/fullstop113/go-web3-demo/model"
	"github.com/fullstop113/go-web3-demo/router"
)

func main() {
	 cfg, err := config.Load()
	 if err != nil {
		log.Fatal(err)
	 }
	model.InitDB()

	r := router.InitRouter()

	log.Println("Server starting on " + cfg.HTTPAddr)
	log.Fatal(r.Run(cfg.HTTPAddr))
}
