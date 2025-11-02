package socket

import "net"

type Client interface {
	Send(command string) (result string, err error)
	Read() (string, error)
	Close() error
}

type client struct {
	conn net.Conn
}

func NewSocketClient(path string) (Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}
	return &client{conn}, nil
}

func (s *client) Send(command string) (result string, err error) {
	_, err = s.conn.Write([]byte(command))
	if err != nil {
		return "", err
	}

	return s.Read()
}

func (s *client) Read() (string, error) {
	buf := make([]byte, 1024)
	n, err := s.conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func (s *client) Close() error {
	return s.conn.Close()
}
