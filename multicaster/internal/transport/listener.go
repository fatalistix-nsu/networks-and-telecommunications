package net

import (
	"encoding/binary"
	"fmt"
	"github.com/fatalistix/slogattr"
	"log/slog"
	"net"
	"strconv"
)

type Handler func(addr string, name string)

type MulticastListener struct {
	log      *slog.Logger
	host     string
	port     int
	handlers map[MessageType]Handler
}

func NewMulticastListener(
	log *slog.Logger,
	host string,
	port int,
) *MulticastListener {
	return &MulticastListener{
		log:      log,
		host:     host,
		port:     port,
		handlers: make(map[MessageType]Handler),
	}
}

func (l *MulticastListener) SetHandler(t MessageType, handler Handler) {
	switch t {
	case JoinOrUpdate:
		l.handlers[JoinOrUpdate] = handler
	case Leave:
		l.handlers[Leave] = handler
	default:
		panic("unexpected message type" + strconv.Itoa(int(t)))
	}
}

func (l *MulticastListener) ListenAndServe() error {
	portStr := strconv.Itoa(l.port)
	hostPort := net.JoinHostPort(l.host, portStr)

	udpAddr, err := net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		return fmt.Errorf("resolv udp address %s: %w", hostPort, err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("listen multicast udp %s: %w", hostPort, err)
	}

	defer l.closeConn(conn)

	buf := make([]byte, typeSize+lengthSize+maxNameSize)

	for {
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("read from %s: %w", hostPort, err)
		}

		msgType := MessageType(buf[0])
		msgLen := binary.BigEndian.Uint16(buf[1:3])
		payload := buf[3 : 3+msgLen]

		h, ok := l.handlers[msgType]
		if !ok {
			return fmt.Errorf("no handler for message type: %s", strconv.Itoa(int(msgType)))
		}

		h(addr.String(), string(payload))
	}
}

func (l *MulticastListener) closeConn(conn *net.UDPConn) {
	err := conn.Close()
	if err != nil {
		l.log.Warn("close connection", slogattr.Err(err))
	}
}

func (l *MulticastListener) Shutdown() {

}
