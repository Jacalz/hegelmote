package remote

import "bytes"

type mockConnection struct {
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

func newControlMock() (*Control, *mockConnection) {
	control := &Control{}
	adapter := &mockConnection{}
	control.conn = adapter
	return control, adapter
}
