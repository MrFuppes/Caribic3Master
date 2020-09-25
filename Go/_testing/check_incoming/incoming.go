package main

import (
	mess "car3-master/message"
	"fmt"
	"log"
	"net"
	"time"
)

func recvMsg(conn *net.UDPConn, ch chan []byte) {
	fmt.Println("starting listener...")
	for {
		buf := make([]byte, 1024)
		n, remoteaddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(time.Now().UTC().Format(isoFmtMilli), "- received message from", remoteaddr)
		ch <- buf[:n]
	}
}

const isoFmtMilli = "2006-01-02T15:04:05.000Z"

// a simple UDP listener that checks if the cRIO sends something and if it is a valid message
func main() {
	// where to expect data
	masAddr, err := net.ResolveUDPAddr("udp", "192.168.1.1:16001")
	if err != nil {
		log.Fatal(err)
	}
	con, err := net.ListenUDP("udp", masAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer con.Close()

	dataIn := make(chan []byte)

	fmt.Println("starting listener on", masAddr)
	go recvMsg(con, dataIn)

	for {
		select {
		case x := <-dataIn:
			fmt.Printf("got %v bytes of data;%v\n", len(x), x)
			parsed, ok := mess.MessageBytes(x).ToMessage()
			if ok {
				fmt.Printf("is valid message;\n%s\n", parsed.String())
			} else {
				fmt.Println("not a valid message according to protocol")
			}
		case <-time.After(5 * time.Second):
			fmt.Println("timeout (5s)")
		}
	}
}
