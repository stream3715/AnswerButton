package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	addr, err := net.ResolveIPAddr("ip4", "internal.kaijudoumei.com")
	if err != nil {
		fmt.Println("Resolve error ", error.Error(err))
		os.Exit(1)
	}
	conn, err := net.Dial("udp", addr.IP.String()+":8804")
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer conn.Close()
	time := time.Now().UnixNano()
	n, err := conn.Write([]byte("lagtest|" + strconv.FormatInt(time, 10)))
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	recvBuf := make([]byte, 1024)

	n, err = conn.Read(recvBuf)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Printf("Received data: %s", string(recvBuf[:n]))
}
