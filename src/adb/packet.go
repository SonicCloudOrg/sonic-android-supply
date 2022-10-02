package adb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
)

/**
 * @link https://github.com/codeskyblue/fa
 */

func swapUint32(n uint32) uint32 {
	var i uint32
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, n)
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

// Packet is a meta for adb connect
type Packet struct {
	Command string
	Arg0    uint32
	Arg1    uint32
	Body    []byte
}

func (pkt Packet) magic() []byte {
	return xorBytes([]byte(pkt.Command), []byte{0xff, 0xff, 0xff, 0xff})
}

func (pkt Packet) checksum() uint32 {
	sum := uint32(0)
	for _, c := range pkt.Body {
		sum += uint32(c)
	}
	return sum
}

func (pkt Packet) length() uint32 {
	return uint32(len(pkt.Body))
}

func (pkt Packet) BodySkipNull() []byte {
	if len(pkt.Body) >= 1 && pkt.Body[len(pkt.Body)-1] == byte(0) {
		return pkt.Body[0 : len(pkt.Body)-1]
	}
	return pkt.Body
}

func (pkt Packet) EncodeToBytes() []byte {
	payload := pkt.Body // append(pkt.Body, byte(0x00))
	buf := bytes.NewBuffer(make([]byte, 0, 24+pkt.length()))
	if len(pkt.Command) != 4 {
		panic("Invalid command " + strconv.Quote(pkt.Command))
	}
	binary.Write(buf, binary.LittleEndian, []byte(pkt.Command))
	binary.Write(buf, binary.LittleEndian, pkt.Arg0)
	binary.Write(buf, binary.LittleEndian, pkt.Arg1)
	binary.Write(buf, binary.LittleEndian, pkt.length())
	binary.Write(buf, binary.LittleEndian, pkt.checksum())
	binary.Write(buf, binary.LittleEndian, pkt.magic())
	buf.Write(payload)
	return buf.Bytes()
}

func (pkt Packet) WriteTo(wr io.Writer) (n int, err error) {
	return wr.Write(pkt.EncodeToBytes())
}

func (pkt Packet) DumpToStdout() {
	fmt.Printf("cmd:%s arg0:%d arg1:%d\n", pkt.Command, pkt.Arg0, pkt.Arg1)
	dumper := hex.Dumper(os.Stdout)
	dumper.Write(pkt.Body)
	dumper.Close()
}
