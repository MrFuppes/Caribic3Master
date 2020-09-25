package main

import (
	inst "car3-master/instrument"
	mess "car3-master/message"
	"car3-master/state"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func recvMsg(conn *net.UDPConn, ch chan []byte) {
	for {
		buf := make([]byte, 128)
		_, remoteaddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(time.Now().UTC().Format(isoFmtMilli), "- received message from", remoteaddr)
		ch <- buf
	}
}

func sendMsg(conn *net.UDPConn, addr *net.UDPAddr, msg []byte) {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Now().UTC().Format(isoFmtMilli), "- sent message to", addr)
}

var recv = flag.String("listen", "192.168.1.64:16064", "listen-on address")
var send = flag.String("reply", "192.168.1.1:16001", "reply-to address")

const isoFmtMilli = "2006-01-02T15:04:05.000Z"

func main() {
	flag.Parse()

	ins := inst.NewInstr()
	ins.Name = "Test-Instrument"
	ins.Address = *recv
	ins.ID = 64
	ins.WUallowed = true
	ins.State = state.Idle
	fmt.Printf("I'm an instrument at (IP:Port) %s\n ...and wait for messages from %s\n", *recv, *send)
	fmt.Printf("my config is\n%v\n", ins)

	// to-do: implement master as an instrument
	masAddr, err := net.ResolveUDPAddr("udp", *send)
	if err != nil {
		log.Fatal(err)
	}

	insAddr, err := ins.ResolveUDPAddr()
	if err != nil {
		log.Fatal(err)
	}

	reply := mess.Message{}
	reply.SendAddr = *insAddr
	if err != nil {
		log.Fatal(err)
	}
	reply.RecvAddr = *masAddr
	reply.MsgType = uint8(4)
	// reply.Data = []byte("hello from instrument!")

	// listen on instrument address
	con, err := net.ListenUDP("udp", insAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()
	dataIn := make(chan []byte)

	// listen for a message from masAddr
	go recvMsg(con, dataIn)

	for {
		select {
		case x := <-dataIn:
			// data on the channel :)
			parsed, _ := mess.MessageBytes(x).ToMessage()
			fmt.Printf(">>> got status SET: %s with timestamp %s\n",
				string(parsed.Data), parsed.Timestamp.Format(isoFmtMilli))
			// reply from insAddr that the message was received
			fmt.Println("sending reply...")
			reply.Timestamp = time.Now().UTC()
			reply.Data = parsed.Data
			sendMsg(con, masAddr, reply.ToBytes())
		case <-time.After(5 * time.Second):
			fmt.Println("timeout 5s")
		}
	}

}
