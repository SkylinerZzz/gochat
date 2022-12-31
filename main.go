package main

import (
	"flag"
	"gochat/common"
	_ "gochat/model"
	"gochat/pkg/adapter"
	"gochat/pkg/task"
	"gochat/router"
	"gochat/util"
	"strconv"
	"time"
)

var port = flag.Int("port", 8080, "port")

func main() {
	// allocate port
	flag.Parse()
	addr := ":" + strconv.Itoa(*port)

	// init server
	util.Init("./config")

	// start queue task
	dispatcher := task.NewDispatcher(util.RedisQueue)
	dispatcherAdapter := adapter.NewQueueTaskAdapter(dispatcher, util.RedisQueue, common.DATABUS_DISPATCHER, 3*time.Second, 100, adapter.NewLogger())
	go dispatcherAdapter.Start()
	contentHandler := task.NewContentHandler(util.RedisQueue)
	contentHandlerAdapter := adapter.NewQueueTaskAdapter(contentHandler, util.RedisQueue, common.DATABUS_CONTENT_HANDLER, 3*time.Second, 100, adapter.NewLogger())
	go contentHandlerAdapter.Start()

	// init router
	r := router.Init()
	r.Run(addr)
}
