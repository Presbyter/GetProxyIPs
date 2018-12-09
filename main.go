package main

import (
	"get_proxy_ips/handler"
	"runtime"
)

func main() {
	go handler.GetIpFromSource()

	runtime.Goexit()
}
