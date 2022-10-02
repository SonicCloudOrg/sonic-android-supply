package adb

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/**
 * @link https://github.com/codeskyblue/fa
 */

type Client struct {
	Addr string
}

func NewClient(addr string) *Client {
	if addr == "" {
		addr = "127.0.0.1:5037"
	}
	return &Client{
		Addr: addr,
	}
}

func (c *Client) dial() (conn *ADBConn, err error) {
	nc, err := net.DialTimeout("tcp", c.Addr, 2*time.Second)
	if err != nil {
		if err = c.StartServer(); err != nil {
			err = errors.Wrap(err, "adb start-server")
			return
		}
		nc, err = net.DialTimeout("tcp", c.Addr, 2*time.Second)
	}
	return NewADBConn(nc), err
}

func (c *Client) roundTrip(data string) (conn *ADBConn, err error) {
	conn, err = c.dial()
	if err != nil {
		return
	}
	if len(data) > 0 {
		err = conn.Encode([]byte(data))
	}
	return
}

func (c *Client) roundTripSingleResponse(data string) (string, error) {
	conn, err := c.roundTrip(data)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if err := conn.CheckOKAY(); err != nil {
		return "", err
	}
	return conn.DecodeString()
}

// ServerVersion returns int. 39 means 1.0.39
func (c *Client) ServerVersion() (v int, err error) {
	verstr, err := c.roundTripSingleResponse("host:version")
	if err != nil {
		return
	}
	_, err = fmt.Sscanf(verstr, "%x", &v)
	return
}

type DeviceState string

const (
	StateUnauthorized = DeviceState("unauthorized")
	StateDisconnected = DeviceState("disconnected")
	StateOffline      = DeviceState("offline")
	StateOnline       = DeviceState("device")
)

// ListDevices returns the list of connected devices
func (c *Client) ListDevices() (devs []*Device, err error) {
	lines, err := c.roundTripSingleResponse("host:devices")
	if err != nil {
		return nil, err
	}

	devs = make([]*Device, 0)
	for _, line := range strings.Split(lines, "\n") {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		devs = append(devs, c.Device(DeviceWithSerial(parts[0])))
	}
	return
}

func (c *Client) StartServer() (err error) {
	cmd := exec.Command("adb", "start-server")
	return cmd.Run()
}

// KillServer tells the server to quit immediately
func (c *Client) KillServer() error {
	conn, err := c.roundTrip("host:kill")
	if err != nil {
		if _, ok := err.(net.Error); ok { // adb is already stopped if connection refused
			return nil
		}
		return err
	}
	defer conn.Close()
	return conn.CheckOKAY()
}

func (c *Client) Device(descriptor DeviceDescriptor) *Device {
	return &Device{
		client:     c,
		descriptor: descriptor,
	}
}

func (c *Client) DeviceWithSerial(serial string) *Device {
	return c.Device(DeviceWithSerial(serial))
}
