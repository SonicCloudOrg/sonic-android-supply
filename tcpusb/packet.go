package tcpusb

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var (
	A_SYNC = uint32(0x434e5953)
	A_CNXN = uint32(0x4e584e43)
	A_OPEN = uint32(0x4e45504f)
	A_OKAY = uint32(0x59414b4f)
	A_CLSE = uint32(0x45534c43)
	A_WRTE = uint32(0x45545257)
	A_AUTH = uint32(0x48545541)
)

type Packet struct {
	check   int
	magic   uint32
	command uint32
	data    *bytes.Buffer
}

func (p *Packet) VerifyChecksum() bool {
	if p.check == 0 {
		return true
	} else {
		return p.check == CheckSun(p.data)
	}
}

func (p *Packet) VerifyMagic() bool {
	return p.magic == Magic(p.command)
}

func (p *Packet) GetType() string {
	switch p.command {
	case A_SYNC:
		return "SYNC"
	case A_CNXN:
		return "CNXN"
	case A_OPEN:
		return "OPEN"
	case A_OKAY:
		return "OKAY"
	case A_CLSE:
		return "CLSE"
	case A_WRTE:
		return "WRTE"
	case A_AUTH:
		return "AUTH"
	}
	return fmt.Sprintf("Unknown command %d", p.command)
}

func CheckSun(data *bytes.Buffer) int {
	var sum int
	if data != nil {
		for _, char := range data.Bytes() {
			sum += int(char)
		}
	}
	return sum
}

func Magic(command uint32) uint32 {
	return (command ^ 0xffffffff) >> 0
}

func Pack(command, arg0, arg1 uint32, data *bytes.Buffer) []byte {
	chunk := new(bytes.Buffer)
	// 1-4 command
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, command)
	chunk.Write(b)
	// 5-8 arg0
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, arg0)
	chunk.Write(b)
	// 9-12 arg1
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, arg1)
	chunk.Write(b)

	if data != nil {
		// 13-16 data size
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(data.Len()))
		chunk.Write(b)
		// 17-20 checksum
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(CheckSun(data)))
		chunk.Write(b)

	} else {
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, 0)
		chunk.Write(b)

		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, 0)
		chunk.Write(b)
	}
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, Magic(command))
	chunk.Write(b)

	return chunk.Bytes()
}
