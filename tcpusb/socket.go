package tcpusb

const (
	UInt32Max        = 0xffffffff
	UInt16Max        = 0xffff
	AuthToken        = 1
	AuthSignature    = 2
	AuthRSAPublicKey = 3
	TokenLength      = 20
)

type Socket struct {
	ended bool
}
