package main

import (
	"fmt"
	"net"
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", net.JoinHostPort("224.0.0.10", "6969"))
	conn, _ := net.ListenMulticastUDP("udp", nil, addr)

	buffer := make([]byte, 1024*1024)

	conn.Read(buffer)

	fmt.Println(string(buffer))
}
