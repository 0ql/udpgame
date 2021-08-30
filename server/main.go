package main

import (
	"fmt"
	"server/netUtils"
	"time"
)

var (
	gameStart time.Time
)

func createServer(port string) {
	gameStart = time.Now()
	fmt.Println("Server is starting...")

	conBundle := netUtils.NewConBundle(10, gameStart)
	go conBundle.ConnectionRemover()
	go conBundle.CreateTCPlistener(":8080")

	udpLn := netUtils.NewUDPListener(":8080", gameStart, &conBundle)
	go udpLn.HandleUDPPackets()
	udpLn.SendUDPStatePackets(10)
}

func main() {
	createServer(":8080")
}
