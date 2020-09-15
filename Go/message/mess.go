// Package mess implements basic types and functionality to handle messages
// on the Container-Lab's LAN.
// Basic package except for general utils from package typeconv.
package mess

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/adler32"
	"net"
	"time"

	"github.com/MrFuppes/Go_General_Use/typeconv"
)

// Message - a struct to hold all the information of a UDP message.
type Message struct {
	SendAddr  net.UDPAddr // who sent
	RecvAddr  net.UDPAddr // who received
	PackLen   uint16
	Timestamp time.Time // the moment the message was sent
	MsgType   uint8     // 0 status set, 1 status is, 2 ambient params, 3 measured data
	Data      []byte
	Checksum  uint32 // Adler32 or CRC32
}

// a variable to hold the checksum function.
var csFunc func([]byte) uint32 = adler32.Checksum // crc32.ChecksumIEEE

// a constant holding the base length of a message (without data).
const msgBaseLen = 27

// ToBytes - method of Message; cast it to a slice of bytes.
// Numbers in big endian byte order.
func (msgStrct *Message) ToBytes() MessageBytes {
	msg := make([]byte, len(msgStrct.Data)+msgBaseLen)
	// sender and receiver addresses:
	copy(msg[0:6], UDPAddrToBytes(msgStrct.SendAddr))
	copy(msg[6:12], UDPAddrToBytes(msgStrct.RecvAddr))
	// length:
	binary.BigEndian.PutUint16(msg[12:14], uint16(len(msgStrct.Data)+msgBaseLen))
	// timestamp:
	t := typeconv.Float64toBytesBE(typeconv.TimetoPOSIX(msgStrct.Timestamp))
	copy(msg[14:22], t)
	// msg type:
	msg[22] = msgStrct.MsgType
	// data; must be byte slice already
	copy(msg[23:23+len(msgStrct.Data)], msgStrct.Data)
	// checksum
	cs := csFunc(msg[:len(msg)-4])
	binary.BigEndian.PutUint32(msg[23+len(msgStrct.Data):msgBaseLen+len(msgStrct.Data)], cs)
	return msg
}

// UDPAddrToBytes - method of Address struct; cast it to a slice of bytes.
// Port is represented as UInt16, big endian.
func UDPAddrToBytes(addrStrct net.UDPAddr) AddressBytes {
	result := make([]byte, 6)
	copy(result[0:4], net.IP.To4(addrStrct.IP))
	binary.BigEndian.PutUint16(result[4:6], uint16(addrStrct.Port))
	return result
}

// AddressBytes - the source for Address struct.
type AddressBytes []byte

// ToAddress - cast AddressBytes (4 bytes IP + 2 bytes port) to Address struct.
func (b AddressBytes) ToAddress() net.UDPAddr {
	var result net.UDPAddr
	result.IP = net.IP(b[0:4])
	result.Port = int(binary.BigEndian.Uint16(b[4:6]))
	return result
}

// MessageBytes - the source for Message struct.
type MessageBytes []byte

// ToMessage - parse message bytes to a message struct.
func (msg MessageBytes) ToMessage() (Message, bool) {
	var result Message
	if len(msg) < msgBaseLen { // assert sufficient length
		return result, false
	}
	l := binary.BigEndian.Uint16(msg[12:14])
	if len(msg) < int(l) { // assert correct length
		return result, false
	}
	if len(msg) > int(l) { // truncate message if too long
		msg = msg[:int(l)]
	}
	cs := csFunc(msg[:len(msg)-4])
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, cs)
	if !bytes.HasSuffix(msg, b) { // assert msg has checksum suffix
		return result, false
	}
	// if not returned by now, everything is fine :)
	result.SendAddr = AddressBytes(msg[0:6]).ToAddress()
	result.RecvAddr = AddressBytes(msg[6:12]).ToAddress()
	result.PackLen = l
	result.Timestamp = typeconv.POSIXtoTime(typeconv.Float64fromBytesBE(msg[14:22]))
	result.MsgType = msg[22]
	result.Data = msg[23 : len(msg)-4]
	result.Checksum = cs
	return result, true
}

// CheckForSig check message bytes if they contain sender - receiver signature.
func (msg MessageBytes) CheckForSig(sendAddr net.UDPAddr, recvAddr net.UDPAddr) (MessageBytes, bool) {
	signature := append(UDPAddrToBytes(sendAddr), UDPAddrToBytes(recvAddr)...)
	idx := bytes.Index(msg, signature)
	if idx == -1 {
		return msg, false
	}
	return msg[idx:], true
}

// String a method to print a readable representation of the message struct.
func (msgStrct *Message) String() string {
	repr := fmt.Sprintf("From Addr:\t%v\n", msgStrct.SendAddr)
	repr += fmt.Sprintf("To Addr:\t%v\n", msgStrct.RecvAddr)
	repr += fmt.Sprintf("Length (bytes):\t%v\n", msgStrct.PackLen)
	repr += fmt.Sprintf("Timestamp:\t%s\n", msgStrct.Timestamp)
	repr += fmt.Sprintf("Type:\t\t%v\n", msgStrct.MsgType)
	repr += fmt.Sprintf("Data (bytes):\t%v\n", msgStrct.Data)
	repr += fmt.Sprintf("Checksum:\t%v\n", msgStrct.Checksum)
	return repr
}
