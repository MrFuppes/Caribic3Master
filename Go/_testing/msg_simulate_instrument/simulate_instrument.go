package main

import (
	"car3-master/msg_parsing/msgutils"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func recvMsg(conn *net.UDPConn, ch chan []byte) {
	fmt.Println("starting listener...")
	for {
		buf := make([]byte, 128)
		_, remoteaddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read a message from %v\n", remoteaddr)
		ch <- buf
	}
}

func sendMsg(conn *net.UDPConn, addr *net.UDPAddr, msg []byte) {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Now().UTC().Format(time.RFC3339), "- sent msg to", addr)
}

var recv = flag.String("listen", "192.168.1.64:16164", "listen-on address")
var send = flag.String("reply", "192.168.1.1:16101", "reply-to address")

func main() {
	flag.Parse()

	masAddr, err := net.ResolveUDPAddr("udp", *send)
	if err != nil {
		log.Fatal(err)
	}
	insAddr, err := net.ResolveUDPAddr("udp", *recv)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(masAddr, " ? ->", insAddr)

	// listen on instrument address
	con, err := net.ListenUDP("udp", insAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer con.Close()
	dataIn := make(chan []byte)

	// listen for a message from masAddr
	go recvMsg(con, dataIn)

	reply := msgutils.Message{}
	reply.SendAddr = *insAddr
	reply.RecvAddr = *masAddr
	reply.MsgType = uint8(4)
	reply.Data = []byte("hello from instrument!")

	for {
		select {
		case x := <-dataIn:
			// data on the channel :)
			parsed, _ := msgutils.MessageBytes(x).ToMessage()
			fmt.Println("got data:", string(parsed.Data))
			// reply from insAddr that the message was received
			fmt.Println("sending reply...")
			reply.Timestamp = time.Now().UTC()
			sendMsg(con, masAddr, reply.ToBytes())
		case <-time.After(5 * time.Second):
			fmt.Println("timeout 5s")
		}
	}

}
