package comm

import (
	"fmt"
	"log"
	"net"
)

type Connection struct {
	LocalPort  string
	RemotePort string
	RemoteIp   string
}

func SimpleServer(connection Connection) {
	fmt.Println("implement connections: ")
	fmt.Println("local: Port")
	fmt.Println(connection.LocalPort)
	fmt.Println("remote client: IP&Port")
	fmt.Println(connection.RemoteIp, ":", connection.RemotePort)
	/* handlling the local result from the real client */
	go func() {
		listen, err := net.Listen("tcp", "0.0.0.0"+":"+connection.LocalPort)
		if err != nil {
			log.Panic("Failed to bind to port 8080 ", err)
		}
		defer listen.Close()

		for {
			connUp, err := listen.Accept()
			if err != nil {
				log.Panic("Error accepting connection: ", err)
			}

			ChanUp := make(chan []byte, 32)
			ChanDown := make(chan []byte, 32)

			go func() {
				connDown, err := net.Dial("tcp", connection.RemoteIp+":"+connection.RemotePort)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				go func() {
					for {
						buf, ok := <-ChanUp
						if !ok {
							connDown.Close()
							return
						}
						connDown.Write(buf)
					}
				}()
				go func() {
					for {
						buf := make([]byte, 1024)
						if n, err := connDown.Read(buf); err == nil {
							ChanDown <- buf[:n]
						} else {
							close(ChanDown)
							return
						}
					}
				}()
			}()
			go func() {
				for {
					buf := make([]byte, 1024)
					if n, err := connUp.Read(buf); err == nil {
						ChanUp <- buf[:n]
					} else {
						close(ChanUp)
						return
					}
				}
			}()
			go func() {
				for {
					buf, ok := <-ChanDown
					if !ok {
						connUp.Close()
						return
					}
					connUp.Write(buf)
				}
			}()
		}
	}()
}
