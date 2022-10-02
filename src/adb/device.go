package adb

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	shellquote "github.com/kballard/go-shellquote"
)

/**
 * @link https://github.com/codeskyblue/fa
 */

// Device
type Device struct {
	descriptor DeviceDescriptor
	client     *Client
}

func (d *Device) String() string {
	return d.descriptor.String()
	// return fmt.Sprintf("%s:%v", ad.serial, ad.State)
}

func (d *Device) Serial() (serial string, err error) {
	return
}

// OpenTransport is a low level function
// Connect to adbd.exe and send <host-prefix>:transport and check OKAY
// conn should be Close after using
func (d *Device) OpenTransport() (conn *ADBConn, err error) {
	req := "host:" + d.descriptor.getTransportDescriptor()
	conn, err = d.client.roundTrip(req)
	if err != nil {
		return
	}
	conn.CheckOKAY()
	if conn.Err() != nil {
		conn.Close()
	}
	return conn, conn.Err()
}

func (d *Device) OpenShell(cmd string) (rwc io.ReadWriteCloser, err error) {
	req := "host:" + d.descriptor.getTransportDescriptor()
	conn, err := d.client.roundTrip(req)
	if err != nil {
		return
	}
	conn.CheckOKAY()
	conn.EncodeString("shell:" + cmd)
	conn.CheckOKAY()
	if conn.Err() != nil {
		conn.Close()
	}
	return conn, conn.Err()
}

func (d *Device) RunCommand(args ...string) (output string, err error) {
	cmd := shellquote.Join(args...)
	rwc, err := d.OpenShell(cmd)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(rwc)
	if err != nil {
		return
	}
	return string(data), err
}

// ServeTCP acts as adbd(Daemon) for adb connect
func (d *Device) ServeTCP(in net.Conn) {
	NewSession(in, d).Serve() // conn will be Closed inside
}

type adbFileInfo struct {
	name  string
	mode  os.FileMode
	size  uint32
	mtime time.Time
}

func (f *adbFileInfo) Name() string {
	return f.name
}

func (f *adbFileInfo) Size() int64 {
	return int64(f.size)
}
func (f *adbFileInfo) Mode() os.FileMode {
	return f.mode
}

func (f *adbFileInfo) ModTime() time.Time {
	return f.mtime
}

func (f *adbFileInfo) IsDir() bool {
	return f.mode.IsDir()
}

func (f *adbFileInfo) Sys() interface{} {
	return nil
}

func (d *Device) Stat(path string) (info os.FileInfo, err error) {
	req := "host:" + d.descriptor.getTransportDescriptor()
	conn, err := d.client.roundTrip(req)
	if err != nil {
		return
	}
	defer conn.Close()
	if err = conn.CheckOKAY(); err != nil {
		return
	}
	conn.EncodeString("sync:")
	conn.CheckOKAY()
	conn.WriteObjects("STAT", uint32(len(path)), path)

	id, err := conn.ReadNString(4)
	if err != nil {
		return
	}
	if id != "STAT" {
		return nil, fmt.Errorf("Invalid status: %q", id)
	}
	adbMode, _ := conn.ReadUint32()
	size, _ := conn.ReadUint32()
	seconds, err := conn.ReadUint32()
	if err != nil {
		return nil, err
	}
	return &adbFileInfo{
		name:  path,
		size:  size,
		mtime: time.Unix(int64(seconds), 0).Local(),
		mode:  fileModeFromAdb(adbMode),
	}, nil
}

type PropValue string

func (p PropValue) Bool() bool {
	return p == "true"
}

var propertyRE = regexp.MustCompile(`\[(.+)\]: \[(.+)\]`)

func (ad *Device) Properties() (props map[string]PropValue, err error) {
	props = make(map[string]PropValue)
	output, err := ad.RunCommand("getprop")
	if err != nil {
		return
	}
	for _, line := range strings.Split(output, "\n") {
		parts := propertyRE.FindStringSubmatch(line)
		if len(parts) != 3 {
			continue
		}
		key, val := parts[1], parts[2]
		props[key] = PropValue(val)
	}
	return
}
