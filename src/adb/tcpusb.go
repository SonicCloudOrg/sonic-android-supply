// Ref link
// https://github.com/openstf/adbkit/blob/master/src/adb/tcpusb/socket.coffee
package adb

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

/**
 * @link https://github.com/codeskyblue/fa
 */

// Session created when adb connected
type Session struct { // adbSession
	device        *Device
	conn          net.Conn
	signature     []byte
	err           error
	token         []byte
	version       uint32
	maxPayload    uint32
	remoteAddress string
	services      map[uint32]*TransportService

	mu             sync.Mutex
	tmpLocalIdLock sync.Mutex
	tmpLocalId     uint32
}

func NewSession(conn net.Conn, device *Device) *Session {
	// generate challenge
	token := make([]byte, TOKEN_LENGTH)
	rand.Read(token)
	log.Println("Create challenge", base64.StdEncoding.EncodeToString(token))

	return &Session{
		device:        device,
		conn:          conn,
		token:         token,
		version:       1,
		remoteAddress: conn.RemoteAddr().String(),
		services:      make(map[uint32]*TransportService),
	}
}

func (s *Session) nextLocalId() uint32 {
	s.tmpLocalIdLock.Lock()
	defer s.tmpLocalIdLock.Unlock()
	s.tmpLocalId += 1
	return s.tmpLocalId
}

func (s *Session) writePacket(cmd string, arg0, arg1 uint32, body []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock() // FIXME(ssx): need to improve performance
	if s.err != nil {
		return s.err
	}
	_, s.err = Packet{
		Command: cmd,
		Arg0:    arg0,
		Arg1:    arg1,
		Body:    body,
	}.WriteTo(s.conn)
	return s.err
}

func (s *Session) Serve() {
	defer s.conn.Close()
	pr := NewPacketReader(s.conn)

	for pkt := range pr.C {
		switch pkt.Command {
		case _CNXN:
			s.onConnection(pkt)
		case _AUTH:
			s.onAuth(pkt)
		case _OPEN:
			s.onOpen(pkt)
		case _OKAY, _WRTE:
			s.forwardServicePacket(pkt)
		case _CLSE:
			s.forwardServicePacket(pkt)
			s.mu.Lock()
			delete(s.services, pkt.Arg1)
			s.mu.Unlock()
		default:
			s.err = errors.New("unknown cmd: " + pkt.Command)
		}
		if s.err != nil {
			log.Printf("unexpect err: %v", s.err)
			break
		}
	}
	log.Println("session closed")
}

func (sess *Session) onConnection(pkt Packet) {
	sess.version = pkt.Arg0
	log.Printf("Version: %x", pkt.Arg0)
	maxPayload := pkt.Arg1
	// log.Println("MaxPayload:", maxPayload)
	if maxPayload > 0xFFFF { // UINT16_MAX
		maxPayload = 0xFFFF
	}
	sess.maxPayload = maxPayload
	// log.Println("MaxPayload:", maxPayload)
	sess.err = sess.writePacket(_AUTH, AUTH_TOKEN, 0, sess.token)
	// pkt.DumpToStdout()
}

func (sess *Session) authVerified() {
	version := swapUint32(1)
	// FIXME(ssx): need device.Properties()
	props, _ := sess.device.Properties()
	connProps := make([]string, 0, 3)
	for _, propName := range []string{
		"ro.product.name",
		"ro.product.model",
		"ro.product.device",
	} {
		connProps = append(connProps, fmt.Sprintf("%s=%s", propName, props[propName]))
	}
	// connProps = append(connProps, "features=cmd") //,stat_v2,shell_v2")
	deviceBanner := "device"
	payload := fmt.Sprintf("%s::%s", deviceBanner, strings.Join(connProps, ";"))
	// id := "device::;;\x00"
	sess.err = sess.writePacket(_CNXN, version, sess.maxPayload, []byte(payload))
	Packet{_CNXN, sess.version, sess.maxPayload, []byte(payload)}.DumpToStdout()
}

func (sess *Session) onAuth(pkt Packet) {
	log.Println("Handle AUTH")
	switch pkt.Arg0 {
	case AUTH_SIGNATURE:
		sess.signature = pkt.Body
		// The real logic is
		// If already have rsa_publickey, then verify signature, send CNXN if passed
		// If no rsa pubkey, then send AUTH to request it
		// Check signature again and send CNXN if passed
		log.Printf("Receive signature: %s", base64.StdEncoding.EncodeToString(pkt.Body))
		// sess.err = sess.writePacket(_AUTH, AUTH_TOKEN, 0, sess.token)
		sess.authVerified()
	case AUTH_RSAPUBLICKEY:
		if sess.signature == nil {
			sess.err = errors.New("Public key sent before signature")
			return
		}
		log.Printf("Receive public key: %s", pkt.Body)
		// TODO(ssx): parse public key from body and verify signature
		// pkt.DumpToStdout()
		log.Println("receive RSA PublicKey")
		// pkt.DumpToStdout()
		// send deviceId
		// time.Sleep(10 * time.Second)
		// sess.err = errors.New("retry")
		// adb 1.0.40 will show "failed to authenticate to x.x.x.x:5555"
		// but actually connected.
		// sess.authVerified()
	default:
		sess.err = fmt.Errorf("unknown authentication method: %d", pkt.Arg0)
	}
}

func (sess *Session) onOpen(pkt Packet) {
	remoteId := pkt.Arg0
	localId := sess.nextLocalId()
	if len(pkt.Body) < 2 {
		sess.err = errors.New("empty service name")
		return // Not throw error ?
	}
	name := string(pkt.BodySkipNull())
	log.Printf("Calling #%s, remoteId: %d, localId: %d\n", name, remoteId, localId)

	service := &TransportService{
		localId:  localId,
		remoteId: remoteId,
		sess:     sess,
	}

	sess.mu.Lock()
	sess.services[localId] = service
	sess.mu.Unlock()

	service.handle(pkt)
	// pkt.DumpToStdout()
}

func (sess *Session) forwardServicePacket(pkt Packet) {
	sess.mu.Lock()
	service, ok := sess.services[pkt.Arg1] // localId
	sess.mu.Unlock()
	if !ok {
		log.Printf("Receive packet of already closed service: localId: %d\n", pkt.Arg1)
		return
	}
	service.handle(pkt)
}

type TransportService struct {
	sess              *Session
	device            *Device
	transport         *ADBConn
	localId, remoteId uint32
	opened            bool
	ended             bool
	once              sync.Once
}

func (t *TransportService) handle(pkt Packet) {
	switch pkt.Command {
	case _OPEN:
		t.handleOpenPacket(pkt)
	case _OKAY:
		// Just ingore
	case _WRTE:
		t.handleWritePacket(pkt)
	case _CLSE:
		t.handleClosePacket(pkt)
	}
}

func (t *TransportService) writeError(message string) {
	t.writePacket(_WRTE, []byte("FAIL"+fmt.Sprintf(
		"%04x%s", len(message), message,
	)))
}

func (t *TransportService) handleOpenPacket(pkt Packet) {
	t.writePacket(_OKAY, nil)

	serviceName := string(pkt.BodySkipNull())
	if strings.HasPrefix(serviceName, "reverse:") {
		failMessage := "reverse service not supported"
		t.writeError(failMessage)
		t.end()
		return
	}

	var err error
	t.transport, err = t.sess.device.OpenTransport()
	if err != nil {
		t.end()
		return
	}
	t.transport.Encode([]byte(serviceName))

	if err := t.transport.CheckOKAY(); err != nil {
		t.writeError(err.Error())
		t.end()
		return
	}

	go func() {
		buf := make([]byte, t.sess.maxPayload)
		for {
			n, err := t.transport.Read(buf)
			if n > 0 {
				t.writePacket(_WRTE, buf[0:n])
			}
			if err != nil {
				t.end()
				break
			}
		}
	}()
}

func (t *TransportService) handleWritePacket(pkt Packet) {
	t.writePacket(_OKAY, nil)
	t.transport.Write(pkt.Body)
}

func (t *TransportService) handleClosePacket(pkt Packet) {
	t.end()
}

func (t *TransportService) end() {
	t.once.Do(func() {
		t.ended = true
		if t.transport != nil {
			t.transport.Close()
		}
		t.writePacket(_CLSE, nil)
	})
}

func (t *TransportService) writePacket(oper string, data []byte) {
	t.sess.writePacket(oper, t.localId, t.remoteId, data)
}
