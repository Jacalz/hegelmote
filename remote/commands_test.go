package remote

import "bytes"

type testConnAdapter struct {
	bytes.Buffer
}

func (t *testConnAdapter) Close() error {
	return nil
}

func newControlTester() (*Control, *testConnAdapter) {
	control := &Control{}
	adapter := &testConnAdapter{}
	control.conn = adapter
	return control, adapter
}
