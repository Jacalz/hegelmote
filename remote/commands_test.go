package remote

import (
	"bytes"
	"net"

	"github.com/Jacalz/hegelmote/device"
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
func (t *mockConnection) Fill() {
	t.readBuf.WriteString("1234567")
}

// FlashToReader flushes written data into the read buffer.
func (t *mockConnection) FlushToReader() {
	t.readBuf.Reset()
	t.readBuf.Write(t.writeBuf.Bytes())
	t.writeBuf.Reset()
}

func newControlMock() (*Control, *mockConnection) {
	control := &Control{model: device.H95}
	adapter := &mockConnection{}
	control.Conn = adapter
	return control, adapter
}
