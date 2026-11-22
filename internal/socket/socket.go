package socket

import (
	"encoding/json"
	"net"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/socket"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/secret"
)

type Client interface {
	Send(command string) (result string, err error)
	SendCommand(command string, args map[string]string) (result string, err error)
	Read() (string, error)
	Close() error
}

type client struct {
	conn   net.Conn
	secret secret.Provider
}

func NewSocketClient(path string, secret secret.Provider) (Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}
	return &client{conn, secret}, nil
}

func (s *client) Send(command string) (result string, err error) {
	secretBytes, err := s.secret.Secret()
	if err != nil {
		return "", err
	}

	args := map[string]string{}
	msg, err := socket.BuildSignedMessage(command, args, secretBytes)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	_, err = s.conn.Write(data)
	if err != nil {
		return "", err
	}

	return s.Read()
}

func (s *client) SendCommand(command string, args map[string]string) (string, error) {
	secretBytes, err := s.secret.Secret()
	if err != nil {
		return "", err
	}

	msg, err := socket.BuildSignedMessage(command, args, secretBytes)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	_, err = s.conn.Write(data)
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
	s.secret.Close()
	return s.conn.Close()
}
