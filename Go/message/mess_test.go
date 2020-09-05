package mess_test

import (
	. "car3-master/Go/message"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestMsgutils(t *testing.T) {
	var ok bool
	var msgStrct Message

	// example
	// Send status SB from 192.168.1.1 (Master Computer) to instrument at 192.168.1.64
	// must be known: master and instrument address
	masAddr, err := net.ResolveUDPAddr("udp", "192.168.1.1:16101") // {IPAddress: net.ParseIP("192.168.1.1"), UDPPort: 16101}
	fmt.Printf("%v, %T (err: %v)\n", masAddr, masAddr, err)
	insAddr, err := net.ResolveUDPAddr("udp", "192.168.1.64:16164") //Address{IPAddress: net.ParseIP("192.168.1.64"), UDPPort: 16164}
	fmt.Printf("%v, %T (err: %v)\n\n", insAddr, insAddr, err)

	// a sample message of bytes with random byte at beginning and end
	// from the master to the instrument telling it to go to standby (SB)
	msgb := MessageBytes{0xff, 0xC0, 0xA8, 0x01, 0x01, 0x3E, 0xE5, 0xC0, 0xA8, 0x01,
		0x40, 0x3F, 0x24, 0x00, 0x1D, 0x41, 0xD7, 0xC6, 0xBB, 0xB3, 0xA4, 0xD7,
		0xA5, 0x00, 0x53, 0x42, 0x87, 0x51, 0x0A, 0xB8, 0xff}

	// check if msg has sender-receiver signature
	msgb, ok = msgb.CheckForSig(*masAddr, *insAddr)
	fmt.Printf("check For Sig:\n %v %v %v\n\n", ok, len(msgb), msgb)

	// parse the message to struct
	msgStrct, ok = msgb.ToMessage()
	fmt.Printf("msg strct:\n %v %v %v\n\n", ok, msgStrct, string(msgStrct.Data))

	//--------------------------------------------------------------------------
	// now the other way around...
	// construct a reply message confirming SB
	response := Message{}
	response.SendAddr = *insAddr
	response.RecvAddr = *masAddr
	response.Timestamp = time.Now().UTC()
	response.MsgType = uint8(1)
	response.Data = []byte{0x06, 0x53, 0x42} // ASCII ack S B
	response.PackLen = uint16(len(response.Data) + 27)

	fmt.Printf("response strct:\n%v\n\n", response)

	msgOut := response.ToBytes()
	fmt.Printf("response bytes:\nlen %v, %v\n\n", msgOut[12:14], msgOut)

	msgStrct, ok = msgOut.ToMessage()
	if ok {
		fmt.Println(ok, msgOut)
		fmt.Println(msgStrct.Data)
	} else {
		fmt.Println("FAILED!")
	}

}
