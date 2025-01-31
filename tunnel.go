package easytcp

import (
	"errors"
	"net"
	"strconv"
	"time"
)

const MaxMessageSize = 1 << 20

type Tunnel struct {
	conn net.Conn
}

type Listener struct {
	listener net.Listener
}

func (t *Tunnel) ReadMessage() ([]byte, error) {
	length := ""
	for {
		char := make([]byte, 1)
		n, err := t.conn.Read(char)
		if err != nil {
			return nil, err
		}
		if n != 1 {
			return nil, errors.New("failed to read one byte from tunnel")
		}
		if char[0] == ';' {
			left, err := strconv.Atoi(length)
			if err != nil {
				return nil, err
			}
			if left <= 0 || left > MaxMessageSize {
				return nil, errors.New("invalid message length or too large")
			}

			var message []byte
			for left > 0 {
				chunk := make([]byte, left)
				n, err := t.conn.Read(chunk)
				if err != nil {
					return nil, err
				}
				if n == 0 {
					return nil, errors.New("failed to read message bytes from tunnel")
				}

				message = append(message, chunk[:n]...)
				left -= n
			}

			return message, nil
		} else {
			length += string(char[0])
			if len(length) > 7 {
				return nil, errors.New("could not read length from invalid message format")
			}
		}
	}
}

func (t *Tunnel) WriteMessage(message []byte) error {
	if len(message) == 0 || len(message) > MaxMessageSize {
		return errors.New("invalid message size")
	}
	length := strconv.Itoa(len(message))
	data := append([]byte(length+";"), message...)
	n, err := t.conn.Write(data)
	if err != nil || n != len(data) {
		return errors.New("failed to write message to tunnel")
	}
	return nil
}

func (t *Tunnel) Close() error {
	return t.conn.Close()
}

func (l *Listener) Accept() (Tunnel, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return Tunnel{}, err
	}
	return Tunnel{conn: conn}, nil
}

func DialTimeout(address string, timeout time.Duration) (Tunnel, error) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return Tunnel{}, err
	}
	return Tunnel{conn: conn}, nil
}

func Listen(address string) (Listener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return Listener{}, err
	}
	return Listener{listener: listener}, nil
}
