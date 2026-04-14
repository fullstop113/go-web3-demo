package main

import (
	"log"

	"github.com/fullstop113/go-web3-demo/model"
	"github.com/fullstop113/go-web3-demo/router"
	"github.com/fullstop113/go-web3-demo/utils"
)

func main() {
	// 初始化JWT
	utils.InitJWT()

	// 初始化数据库
	model.InitDB()

	// 初始化路由
	r := router.InitRouter()

	log.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}
