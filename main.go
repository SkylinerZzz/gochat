package main

import (
	"github.com/gin-gonic/gin"
	_ "gochat/model"
	"gochat/router"
)

func main() {
	r := gin.Default()
	r = router.InitRouter(r)
	r.Run(":8080")
}
