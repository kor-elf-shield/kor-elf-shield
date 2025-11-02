package socket

import "net"

type Connect interface {
	Read() (string, error)
	Write(command string) error
	Close() error
}

type connect struct {
	conn net.Conn
}

func NewConnect(conn net.Conn) Connect {
	return &connect{
		conn: conn,
	}
}

func (c *connect) Read() (string, error) {
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func (c *connect) Write(command string) error {
	_, err := c.conn.Write([]byte(command))
	return err
}

func (c *connect) Close() error {
	return c.conn.Close()
}
