package adb

import "io"

/**
 * @link https://github.com/codeskyblue/fa
 */

type PacketWriter struct {
	wr io.Writer
}

func NewPacketWriter(w io.Writer) *PacketWriter {
	return &PacketWriter{
		wr: w,
	}
}

func (p *PacketWriter) WritePacket(pkt Packet) error {
	_, err := p.wr.Write(pkt.EncodeToBytes())
	return err
}
