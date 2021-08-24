package main

import (
	"fmt"
	"server/TCPutils"
	"server/UDPutils"
	"time"
)

var (
	gameStart time.Time
)

func createServer(port string) {
	gameStart = time.Now()
	fmt.Println("Server is starting...")

	tcpConBundle := TCPutils.NewTCPConBundle(10, gameStart)
	go tcpConBundle.ConnectionRemover()
	go tcpConBundle.CreateTCPlistener(":8080")

	udpLn := UDPutils.NewUDPListener(":8080", gameStart, &tcpConBundle)
	go udpLn.HandleUDPPackets()
	udpLn.SendUDPStatePackets(3)
}

func main() {
	createServer(":8080")
}
