package main

import (
	inst "car3-master/Go/instrument"
	mess "car3-master/Go/message"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"time"
)

var payload inst.Payload

func init() {
	wd, _ := os.Getwd()
	src := path.Join(wd, "instr_cfg_v20200907.yml")
	payload, _ = inst.PayloadFromYAML(src)
}

func sendMsg(conn *net.UDPConn, addr *net.UDPAddr, msg []byte) {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Now().UTC().Format(time.RFC3339), "- sent msg to", addr)
}

func recvMsg(conn *net.UDPConn, ch chan []byte) {
	fmt.Println("starting listener...")
	for {
		buf := make([]byte, 1024)
		_, remoteaddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read a message from %v \n", remoteaddr)
		ch <- buf
	}
}

var send = flag.String("from", "192.168.1.1:16101", "sender address")
var recv = flag.String("to", "192.168.1.64:16164", "receiver address")
var msg = flag.String("message", "hello from master :)", "the message to send")

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
	fmt.Println(masAddr, " ! ->", insAddr)

	for k, v := range payload {
		fmt.Println(k, v)
	}

	// construct the message to send
	message := mess.Message{}
	message.SendAddr = *masAddr
	message.RecvAddr = *insAddr
	message.MsgType = uint8(3)
	message.Data = []byte(*msg)
	fmt.Printf("message strct:\n%v\n\n", message)

	con, err := net.ListenUDP("udp", masAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer con.Close()
	dataIn := make(chan []byte)

	// simultaneously wait for reply
	go recvMsg(con, dataIn)

	// now keep sending this to insAddr
	fmt.Println("starting sender...")
	for {
		select {
		case x := <-dataIn:
			// data on the channel :)
			parsed, _ := mess.MessageBytes(x).ToMessage()
			fmt.Println("got data:", string(parsed.Data))
		case <-time.After(5 * time.Second):
			fmt.Println("sending message...")
			message.Timestamp = time.Now().UTC()
			sendMsg(con, insAddr, message.ToBytes())
		}
	}

}
