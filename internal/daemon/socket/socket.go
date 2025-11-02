package socket

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type HandleCommand func(command string) error

type Socket interface {
	EnsureNoOtherProcess() error
	Create() error
	Close() error
	Run(ctx context.Context, handleCommand HandleCommand)
}

type socket struct {
	path   string
	logger log.Logger

	listen net.Listener
}

func New(path string, logger log.Logger) (Socket, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	return &socket{
		path:   path,
		logger: logger,
	}, nil
}

func (s *socket) EnsureNoOtherProcess() error {
	if s.path == "" {
		return errors.New("socket file not specified")
	}
	fi, err := os.Lstat(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			// File isn't found, so everything is fine.
			return nil
		}
		return err
	}

	if fi.Mode()&os.ModeSocket == 0 {
		return fmt.Errorf("%s exists but is not a socket", s.path)
	}

	if canConnect(s.path) {
		return errors.New("another process is already listening on the socket")
	}

	if err := os.Remove(s.path); err != nil {
		return err
	}
	s.logger.Warn(fmt.Sprintf("An obsolete socket file %s was found: file removed", s.path))
	return nil
}

func (s *socket) Create() error {
	ln, err := net.Listen("unix", s.path)
	if err != nil {
		return err
	}
	s.listen = ln
	return nil
}

func (s *socket) Close() error {
	if s.path == "" {
		s.logger.Warn("socket file not specified")
		return nil
	}

	s.logger.Debug("Remove socket file")
	if err := s.listen.Close(); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to close socket: %s", err))
		return err
	}

	return nil
}

func (s *socket) Run(ctx context.Context, handleCommand HandleCommand) {
	for {
		conn, err := s.listen.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				// shutdown, exit the acceptance cycle
				s.logger.Debug("The socket has closed. Shutdown")
				return
			default:
				s.logger.Error(fmt.Sprintf("Failed to accept connection: %s", err))
				continue
			}
		}
		go s.handleConn(conn, handleCommand)
	}
}

func (s *socket) handleConn(conn net.Conn, handleCommand HandleCommand) {
	defer func() {
		_ = conn.Close()
	}()
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	cmd := string(buf[:n])

	if err := handleCommand(cmd); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to handle command: %s", err))
	}
}

func canConnect(path string) bool {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return false
	}
	defer func() {
		_ = conn.Close()
	}()

	return true
}
