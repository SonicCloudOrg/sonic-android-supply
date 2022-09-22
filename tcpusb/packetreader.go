package tcpusb

import "io"

type PacketReader struct {
	inBody bool
	buffer io.ByteReader
}
