package main

import (
	"flag"
	"fmt"
	"gochat/common"
	_ "gochat/model"
	"gochat/pkg/adapter"
	"gochat/pkg/task"
	"gochat/router"
	"gochat/util"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
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
	dispatcherAdapter := adapter.NewQueueTaskAdapter(dispatcher, util.RedisQueue, common.DATABUS_DISPATCHER, 3*time.Second, 100, adapter.NewLogger("dispatcher"))
	go dispatcherAdapter.Start()
	contentHandler := task.NewContentHandler(util.RedisQueue)
	contentHandlerAdapter := adapter.NewQueueTaskAdapter(contentHandler, util.RedisQueue, common.DATABUS_CONTENT_HANDLER, 3*time.Second, 100, adapter.NewLogger("content handler"))
	go contentHandlerAdapter.Start()

	// init router
	r := router.Init()
	go r.Run(addr)

	// start pprof
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			fmt.Println("exit pprof")
		}
	}()

	// listen terminate signal
	wg := sync.WaitGroup{}
	wg.Add(1)
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)
	go func() {
		<-exit
		dispatcherAdapter.Terminate()
		contentHandlerAdapter.Terminate()
		wg.Done()
		fmt.Println("exit server")
	}()
	wg.Wait()
}
