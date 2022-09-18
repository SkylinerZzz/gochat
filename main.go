package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "gochat/model"
	"gochat/router"
	"strconv"
)

var port = flag.Int("port", 8080, "port")

func main() {
	flag.Parse()
	addr := ":" + strconv.Itoa(*port)
	fmt.Println(addr)
	r := gin.Default()
	r = router.InitRouter(r)
	r.Run(addr)
}
