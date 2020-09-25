package main

import (
	inst "car3-master/instrument"
	mess "car3-master/message"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v2"
)

var payload inst.Payload

func init() {
	wd, _ := os.Getwd()
	src := path.Join(wd, "instr_cfg_v20200914.yml")
	payload, _ = inst.PayloadFromCfg(src, yaml.Unmarshal)
}

func sendMsg(conn *net.UDPConn, addr *net.UDPAddr, msg []byte) {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Now().UTC().Format(isoFmtMilli), "- sent msg to", addr)
}

func recvMsg(conn *net.UDPConn, ch chan []byte) {
	fmt.Println("starting listener...")
	for {
		buf := make([]byte, 1024)
		_, remoteaddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(time.Now().UTC().Format(isoFmtMilli), "- received message from", remoteaddr)
		ch <- buf
	}
}

var send = flag.String("from", "192.168.1.1:16001", "sender address")
var recv = flag.String("to", "192.168.1.64:16064", "receiver address")
var msg = flag.String("message", "SB", "the message to send")

const isoFmtMilli = "2006-01-02T15:04:05.000Z"

func main() {
	flag.Parse()
	fmt.Printf("I'm the Master (IP:Port) %s\n", *send)

	masAddr, err := net.ResolveUDPAddr("udp", *send)
	if err != nil {
		log.Fatal(err)
	}
	insAddr, err := net.ResolveUDPAddr("udp", *recv)
	if err != nil {
		log.Fatal(err)
	}

	// for k, v := range payload {
	// 	fmt.Println(k, v)
	// }

	fmt.Println(payload[1])

	// construct the message to send
	message := mess.Message{}
	message.SendAddr = *masAddr
	message.RecvAddr = *insAddr
	message.MsgType = uint8(0)
	message.Data = []byte(*msg)
	fmt.Printf("prepared message:\n%s\n\n", message.String())

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
			fmt.Printf(">>> got status IS: %s with timestamp %s\n",
				string(parsed.Data), parsed.Timestamp.Format(isoFmtMilli))
		case <-time.After(5 * time.Second):
			message.Timestamp = time.Now().UTC()
			sendMsg(con, insAddr, message.ToBytes())
		}
	}

}
