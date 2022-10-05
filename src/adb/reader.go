package adb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
)

/**
 * @link https://github.com/codeskyblue/fa
 */

const (
	_SYNC = "SYNC"
	_CNXN = "CNXN"
	_OPEN = "OPEN"
	_OKAY = "OKAY"
	_CLSE = "CLSE"
	_WRTE = "WRTE"
	_AUTH = "AUTH"
	_FAIL = "FAIL"

	UINT16_MAX = 0xFFFF
	UINT32_MAX = 0xFFFFFFFF

	AUTH_TOKEN        = 1
	AUTH_SIGNATURE    = 2
	AUTH_RSAPUBLICKEY = 3

	TOKEN_LENGTH = 20
)

var (
	ErrChecksum   = errors.New("adb: checksum error")
	ErrCheckMagic = errors.New("adb: magic error")
)

func calculateChecksum(data []byte) uint32 {
	sum := uint32(0)
	for _, c := range data {
		sum += uint32(c)
	}
	return sum
}

func xorBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		panic(fmt.Sprintf("xorBytes a:%x b:%x have different size", a, b))
	}
	dst := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst
}

type PacketReader struct {
	C      chan Packet
	reader io.Reader
	err    error
}

func NewPacketReader(reader io.Reader) *PacketReader {
	pr := &PacketReader{
		C:      make(chan Packet, 1),
		reader: reader,
	}
	go pr.drain()
	return pr
}

func (p *PacketReader) Err() error {
	return p.err
}

type errReader struct{}

func (e errReader) Read(p []byte) (int, error) {
	return 0, errors.New("package already read error")
}

func (p *PacketReader) r() io.Reader {
	if p.err != nil { // use p.err to short error checks
		return errReader{}
	}
	return p.reader
}

func (p *PacketReader) drain() {
	defer close(p.C)
	for {
		pkt, err := p.readPacket()
		if err != nil {
			break
		}
		if p.Err() != nil {
			log.Println("packet read error", p.Err())
			break
		}
		p.C <- pkt
	}
}

// Receive packet example
// 00000000  43 4e 58 4e 01 00 00 01  00 00 10 00 23 00 00 00  |CNXN........#...|
// 00000010  3c 0d 00 00 bc b1 a7 b1  68 6f 73 74 3a 3a 66 65  |<.......host::fe|
// 00000020  61 74 75 72 65 73 3d 63  6d 64 2c 73 74 61 74 5f  |atures=cmd,stat_|
// 00000030  76 32 2c 73 68 65 6c 6c  5f 76 32                 |v2,shell_v2|
func (p *PacketReader) readPacket() (pkt Packet, err error) {
	pkt = Packet{
		Command: p.readStringN(4),
		Arg0:    p.readUint32(),
		Arg1:    p.readUint32(),
	}

	var (
		length   = p.readUint32()
		checksum = p.readUint32()
		magic    = p.readN(4)
	)

	pkt.Body = p.readN(int(length))

	if p.err != nil {
		return
	}
	if !bytes.Equal(xorBytes([]byte(pkt.Command), magic), []byte{0xff, 0xff, 0xff, 0xff}) {
		p.err = ErrCheckMagic
		log.Printf("%x %x %x", []byte(pkt.Command), magic, xorBytes([]byte(pkt.Command), magic))
		return
	}
	// log.Printf("cmd:%s, arg0:%x, arg1:%x, len:%d, check:%x, magic:%x",
	// 	pkt.Command, pkt.Arg0, pkt.Arg1, length, checksum, magic)
	if calculateChecksum(pkt.Body) != checksum {
		p.err = ErrChecksum
	}
	return pkt, p.err
}

func (p *PacketReader) readN(n int) []byte {
	buf := make([]byte, n)
	_, p.err = io.ReadFull(p.r(), buf)
	return buf
}

func (p *PacketReader) readStringN(n int) string {
	return string(p.readN(n))
}

func (p *PacketReader) readInt32() int32 {
	var i int32
	p.err = binary.Read(p.r(), binary.LittleEndian, &i)
	return i
}

func (p *PacketReader) readUint32() uint32 {
	var i uint32
	p.err = binary.Read(p.r(), binary.LittleEndian, &i)
	return i
}
