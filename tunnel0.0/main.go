package main

import (
	. "tunnel/comm"
)

func main() {
	var conn Connection
	conn.LocalPort = "8081"
	conn.RemoteIp = "44.114.12.43"
	conn.RemotePort = "8081"
	SimpleServer(conn)
	<-make(chan []byte)
}
