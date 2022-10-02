package adb

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

/**
 * @link https://github.com/codeskyblue/fa
 */

type ADBConn struct {
	rw io.ReadWriter
	io.Closer
	err error
}

func NewADBConn(conn net.Conn) *ADBConn {
	proxyRW := debugProxyConn{
		R:     bufio.NewReader(conn),
		W:     conn,
		Debug: false}

	return &ADBConn{
		rw:     proxyRW,
		Closer: conn,
	}
}

func (conn *ADBConn) Err() error {
	return conn.err
}

func (conn *ADBConn) Read(p []byte) (n int, err error) {
	if conn.err != nil {
		return 0, conn.err
	}
	n, err = conn.rw.Read(p)
	conn.err = err
	return
}

func (conn *ADBConn) Write(p []byte) (n int, err error) {
	if conn.err != nil {
		return 0, conn.err
	}
	n, err = conn.rw.Write(p)
	conn.err = err
	return
}

func (conn *ADBConn) Encode(v []byte) error {
	val := string(v)
	return conn.EncodeString(val)
}

func (conn *ADBConn) EncodeString(s string) error {
	data := fmt.Sprintf("%04x%s", len(s), s)
	_, err := conn.Write([]byte(data))
	return err
}

// write data with little endian
func (conn *ADBConn) WriteLE(v interface{}) error {
	return binary.Write(conn, binary.LittleEndian, v)
}

func (conn *ADBConn) WriteString(s string) (int, error) {
	return conn.Write([]byte(s))
}

// WriteObjects according to type
func (conn *ADBConn) WriteObjects(objs ...interface{}) error {
	var err error
	for _, obj := range objs {
		switch obj.(type) {
		case string:
			_, err = conn.WriteString(obj.(string))
		case uint32, int32, uint16, int16:
			err = conn.WriteLE(obj)
		default:
			err = fmt.Errorf("Unsupported type: %t", obj)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (conn *ADBConn) ReadUint32() (i uint32, err error) {
	err = binary.Read(conn, binary.LittleEndian, &i)
	return
}

func (conn *ADBConn) ReadN(n int) (data []byte, err error) {
	buf := make([]byte, n)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return
	}
	return buf, nil
}

func (conn *ADBConn) ReadNString(n int) (data string, err error) {
	bdata, err := conn.ReadN(n)
	return string(bdata), err
}

func (conn *ADBConn) DecodeString() (string, error) {
	hexlen, err := conn.ReadNString(4)
	if err != nil {
		return "", err
	}
	var length int
	_, err = fmt.Sscanf(hexlen, "%04x", &length)
	if err != nil {
		return "", err
	}
	return conn.ReadNString(length)
}

// CheckOKAY check OKAY, or FAIL
func (conn *ADBConn) CheckOKAY() error {
	status, _ := conn.ReadNString(4)
	switch status {
	case _OKAY:
		return nil
	case _FAIL:
		data, err := conn.DecodeString()
		if err != nil {
			return err
		}
		return errors.Wrap(errors.New(data), "respCheck")
	default:
		return fmt.Errorf("Unexpected response: %s, should be OKAY or FAIL", strconv.Quote(status))
	}
}

type debugProxyConn struct {
	R     io.Reader
	W     io.Writer
	Debug bool
}

func (px debugProxyConn) Write(data []byte) (int, error) {
	if px.Debug {
		m := regexp.MustCompile(`^[-:/0-9a-zA-Z ]+$`)
		if m.Match(data) {
			fmt.Printf("-> %q\n", string(data))
		} else {
			fmt.Printf("-> \\x%x\n", reverseBytes(data))
		}
	}
	return px.W.Write(data)
}

func reverseBytes(b []byte) []byte {
	out := make([]byte, len(b))
	for i, c := range b {
		out[len(b)-i-1] = c
	}
	return out
}

func (px debugProxyConn) Read(data []byte) (int, error) {
	n, err := px.R.Read(data)
	if px.Debug {
		m := regexp.MustCompile(`^[-:/0-9a-zA-Z ]+$`)
		if m.Match(data[0:n]) {
			fmt.Printf("<---- %q\n", string(data[0:n]))
		} else {
			fmt.Printf("<---- \\x%x\n", reverseBytes(data[0:n]))
		}
	}
	return n, err
}
