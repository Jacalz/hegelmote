package remote

import (
	"bytes"
	"net"
	"testing"

	"github.com/Jacalz/hegelmote/device"
	"github.com/alecthomas/assert/v2"
)

type mockConnection struct {
	net.Conn
	readBuf  bytes.Buffer
	writeBuf bytes.Buffer
}

// Read does a read from the read buffer.
func (t *mockConnection) Read(buf []byte) (int, error) {
	return t.readBuf.Read(buf)
}

// Write passes the buffer contents down into the write buffer.
func (t *mockConnection) Write(buf []byte) (int, error) {
	return t.writeBuf.Write(buf)
}

// Close is the same as calling Reset() on both buffers.
func (t *mockConnection) Close() error {
	t.readBuf.Reset()
	t.writeBuf.Reset()
	return nil
}

// Fill fills the read buffer to avoid EOF errors.
func (t *mockConnection) Fill(response string) {
	t.readBuf.Reset()
	t.writeBuf.Reset()
	t.readBuf.WriteString(response)
}

// FlashToReader flushes written data into the read buffer.
func (t *mockConnection) FlushToReader() {
	t.readBuf.Reset()
	t.readBuf.Write(t.writeBuf.Bytes())
	t.writeBuf.Reset()
}

func newControlMock() (*Control, *mockConnection) {
	control := &Control{deviceType: device.H95}
	adapter := &mockConnection{}
	control.conn = adapter
	return control, adapter
}

var parseUint8FromBufTestcases = []struct {
	in    string
	val   uint8
	error bool
}{
	{"-v.0\r", 0, false},
	{"-v.10\r", 10, false},
	{"-v.100\r", 100, false},
	{"-v.255\r", 255, false},
	{"-v.256\r", 0, true},
	{"-v.a23\r", 0, true},
	{"-v.1bc\r", 0, true},
	{"-v.12c\r", 0, true},
	{"-v.100\rxxxxxxxx", 100, false},
}

func TestParseUint8FromBuf(t *testing.T) {
	for _, tc := range parseUint8FromBufTestcases {
		out, err := parseUint8FromBuf([]byte(tc.in))
		assert.Equal(t, tc.val, out)
		assert.Equal(t, tc.error, err != nil)
	}
}
