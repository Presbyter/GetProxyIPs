package main

import (
	"context"
	"get_proxy_ips/common"
	"get_proxy_ips/handler"
	"runtime"
)

func main() {
	common.GetAppBus()

	go handler.GetIpFromSource()
	go handler.CleanIps(context.TODO())
	//
	//runtime.Goexit()
	//sigChan := make(chan os.Signal)
	//signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	//log.Printf("signal received signal %v", <-sigChan)

	//ch := make(chan int, 10)
	//
	//go func() {
	//	for i := 0; i < 100; i++ {
	//		ch <- i
	//		time.Sleep(200 * time.Millisecond)
	//	}
	//}()
	//
	//go func() {
	//	for item := range ch {
	//		log.Println(item)
	//	}
	//}()

	runtime.Goexit()
}
