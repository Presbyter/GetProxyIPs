package main

import (
	"get_proxy_ips/common"
	"get_proxy_ips/handler"
	"runtime"
)

func main() {
	common.GetAppBus()

	go handler.GetIpFromSource()

	runtime.Goexit()
}
