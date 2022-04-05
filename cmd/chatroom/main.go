package main

import (
	// 初始化 global，详细理解 go的初始化顺序
	_ "chatroom/global"
	"chatroom/routes"
	"fmt"
	"log"
	"net/http"
)

var (
	addr   = ":13140"
	outStr = `
    ____               		     _______  
   |    	|    |      /\     	|
   |    	|____|     /  \    	| 
   |    	|    | 	  /----\   	|
   |____	|    |	 /      \  	|


诸葛青的编程之旅：ChatRoom，listening on：%s

`
)

func main() {
	fmt.Printf(outStr, addr)
	routes.RegisterHandle()
	log.Fatal(http.ListenAndServe(addr, nil))
}
