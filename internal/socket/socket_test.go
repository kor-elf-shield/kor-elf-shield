package socket

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/secret_mock"
)

type mockConn struct {
	readData  []byte
	readErr   error
	writeErr  error
	closeErr  error
	written   bytes.Buffer
	readCalls int
}

func (m *mockConn) Read(b []byte) (int, error) {
	m.readCalls++
	if m.readErr != nil {
		return 0, m.readErr
	}
	if len(m.readData) == 0 {
		return 0, io.EOF
	}
	n := copy(b, m.readData)
	m.readData = m.readData[n:]
	return n, nil
}

func (m *mockConn) Write(b []byte) (int, error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return m.written.Write(b)
}

func (m *mockConn) Close() error { return m.closeErr }

func (m *mockConn) LocalAddr() net.Addr                { return &net.UnixAddr{Name: "local", Net: "unix"} }
func (m *mockConn) RemoteAddr() net.Addr               { return &net.UnixAddr{Name: "remote", Net: "unix"} }
func (m *mockConn) SetDeadline(_ time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(_ time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(_ time.Time) error { return nil }

func TestNewSocketClient_Success(t *testing.T) {
	dir := t.TempDir()
	socketPath := filepath.Join(dir, "test.sock")

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("net.Listen failed to create socket: %v", err)
	}
	defer func() {
		_ = ln.Close()
		_ = os.Remove(socketPath)
	}()

	secret := secret_mock.New()
	c, err := NewSocketClient(socketPath, secret)
	if err != nil {
		t.Fatalf("NewSocketClient: failed to create socket client: %v", err)
	}

	if c == nil {
		t.Fatalf("NewSocketClient: socket client is nil")
	}

	_ = c.Close()
}

func TestNewSocketClient_Error(t *testing.T) {
	dir := t.TempDir()
	secret := secret_mock.New()

	socketPath := filepath.Join(dir, "non_existing_test.sock")
	if _, err := NewSocketClient(socketPath, secret); err == nil {
		t.Fatalf("NewSocketClient: expected error for non-existing unix socket path")
	}

	if _, err := NewSocketClient("", secret); err == nil {
		t.Fatalf("NewSocketClient: expected error for invalid unix socket path")
	}
}

func TestClient_Send_Success(t *testing.T) {
	mc := &mockConn{readData: []byte("ok")}
	secret := secret_mock.New()
	c := &client{conn: mc, secret: secret}

	got, err := c.Send("ping")
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if got != "ok" {
		t.Fatalf("unexpected send result: got=%q want=%q", got, "ok")
	}
}

func TestClient_Send_Error(t *testing.T) {
	expErr := errors.New("write failed")
	mc := &mockConn{writeErr: expErr}
	secret := secret_mock.New()
	c := &client{conn: mc, secret: secret}

	_, err := c.Send("ping")
	if !errors.Is(err, expErr) {
		t.Fatalf("expected write error, got: %v", err)
	}
	if mc.readCalls != 0 {
		t.Fatalf("Read must not be called after write error, got readCalls=%d", mc.readCalls)
	}
}

func TestClient_SendCommand_Success(t *testing.T) {
	mc := &mockConn{readData: []byte("ok")}
	secret := secret_mock.New()
	c := &client{conn: mc, secret: secret}

	if got, err := c.SendCommand("ping", make(map[string]string)); err != nil {
		t.Fatalf("Send returned error: %v", err)
	} else if got != "ok" {
		t.Fatalf("unexpected send result: got=%q want=%q", got, "ok")
	}

	mc = &mockConn{readData: []byte("ok")}
	secret = secret_mock.New()
	c = &client{conn: mc, secret: secret}
	if got, err := c.SendCommand("ping", map[string]string{"name": "value"}); err != nil {
		t.Fatalf("Send returned error: %v", err)
	} else if got != "ok" {
		t.Fatalf("unexpected send result: got=%q want=%q", got, "ok")
	}

	mc = &mockConn{readData: []byte("ok")}
	secret = secret_mock.New()
	c = &client{conn: mc, secret: secret}
	if got, err := c.SendCommand("ping", map[string]string{"name": ""}); err != nil {
		t.Fatalf("Send returned error: %v", err)
	} else if got != "ok" {
		t.Fatalf("unexpected send result: got=%q want=%q", got, "ok")
	}
}

func TestClient_SendCommand_Error(t *testing.T) {
	expErr := errors.New("write failed")
	mc := &mockConn{writeErr: expErr}
	secret := secret_mock.New()
	c := &client{conn: mc, secret: secret}

	if _, err := c.SendCommand("ping", make(map[string]string)); !errors.Is(err, expErr) {
		t.Fatalf("expected write error, got: %v", err)
	}
	if mc.readCalls != 0 {
		t.Fatalf("Read must not be called after write error, got readCalls=%d", mc.readCalls)
	}

	mc = &mockConn{writeErr: expErr}
	secret = secret_mock.New()
	c = &client{conn: mc, secret: secret}
	if _, err := c.SendCommand("ping", map[string]string{"name": "value"}); !errors.Is(err, expErr) {
		t.Fatalf("expected write error, got: %v", err)
	}
	if mc.readCalls != 0 {
		t.Fatalf("Read must not be called after write error, got readCalls=%d", mc.readCalls)
	}

	mc = &mockConn{writeErr: expErr}
	secret = secret_mock.New()
	c = &client{conn: mc, secret: secret}
	if _, err := c.SendCommand("ping", map[string]string{"name": ""}); !errors.Is(err, expErr) {
		t.Fatalf("expected write error, got: %v", err)
	}
	if mc.readCalls != 0 {
		t.Fatalf("Read must not be called after write error, got readCalls=%d", mc.readCalls)
	}
}

func TestClient_Read_Success(t *testing.T) {
	mc := &mockConn{readData: []byte("hello")}
	secret := secret_mock.New()
	c := &client{conn: mc, secret: secret}

	got, err := c.Read()
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("unexpected read data: got=%q want=%q", got, "hello")
	}
}

func TestClient_Read_Error(t *testing.T) {
	expErr := errors.New("read failed")
	secret := secret_mock.New()
	mc := &mockConn{readErr: expErr}
	c := &client{conn: mc, secret: secret}

	_, err := c.Read()
	if !errors.Is(err, expErr) {
		t.Fatalf("expected read error, got: %v", err)
	}
}

func TestClient_Close_Success(t *testing.T) {
	secret := secret_mock.New()
	mc := &mockConn{}
	c := &client{conn: mc, secret: secret}

	err := c.Close()
	if err != nil {
		t.Fatalf("expected close error, got: %v", err)
	}
}

func TestClient_Close_Error(t *testing.T) {
	secret := secret_mock.New()

	expErr := errors.New("close")
	mc := &mockConn{closeErr: expErr}
	c := &client{conn: mc, secret: secret}

	err := c.Close()
	if !errors.Is(err, expErr) {
		t.Fatalf("expected close error, got: %v", err)
	}
}
