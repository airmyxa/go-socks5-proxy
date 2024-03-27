package proto

import (
	"fmt"
	"net"
)

type SessionError struct {
	message string
}

func (se SessionError) Error() string {
	return se.message
}

type SessionPool struct {
}

func MakeSessionPool() SessionPool {
	return SessionPool{}
}

func (sp *SessionPool) NewSession(lstream net.Conn) Session {
	return Session{lstream: lstream, rstream: nil}
}

type Session struct {
	lstream net.Conn
	rstream net.Conn
}

func (s *Session) Start(lstream net.Conn) error {
	s.lstream = lstream
	err := s.handshake()
	if err != nil {
		return SessionError{message: fmt.Sprintf("Could not process handshake: %s", err.Error())}
	}

	err = s.pipe()
	if err != nil {
		return SessionError{message: fmt.Sprintf("Error while piping data: %s", err.Error())}
	}

	err = s.shutdown()
	if err != nil {
		return SessionError{message: fmt.Sprintf("Error while trying to shutdown session: %s", err.Error())}
	}

	return nil
}

func (s *Session) handshake() error {
	bytes := make([]byte, 256)

	_, err := s.lstream.Read(bytes)
	if err != nil {
		return SessionError{message: fmt.Sprintf("Client error: %s", err.Error())}
	}

	var csm ClientSuggestedMethods
	err = csm.FromBytes(bytes)
	if err != nil {
		return SessionError{message: fmt.Sprintf("Client error: %s", err.Error())}
	}

	ssm := ServerSelectedMethod{Ver: Ver5, Method: 0}
	size := copy(bytes, ssm.ToBytes())
	bytes = bytes[:size]
	_, err = s.lstream.Write(bytes)
	if err != nil {
		return SessionError{message: fmt.Sprintf("Client error: %s", err.Error())}
	}

	_, err = s.lstream.Read(bytes)
	if err != nil {
		return err
	}

	var csr ClientSocksRequest
	err = csr.FromBytes(bytes)
	if err != nil {
		return err
	}

	addr, err := func() (net.IP, error) {
		switch csr.AddrType {
		case Ipv4:
			return net.IPv4(csr.DestAddr[0], csr.DestAddr[1], csr.DestAddr[2], csr.DestAddr[3]), nil
		case Ipv6:
			return net.IP(csr.DestAddr), nil
		case Domain:
			addr, err := net.ResolveIPAddr("ip", string(csr.DestAddr))
			return addr.IP, err
		}
		return nil, SessionError{message: "Cannot parse ip addr"}
	}()

	if err != nil {
		return err
	}

	s.rstream, err = net.DialTCP("ip", nil, &net.TCPAddr{Zone: "", IP: addr, Port: int(csr.DestPort)})
	if err != nil {
		return err
	}

	ssr := ServerSocksResponse{Ver: Ver5, ReplyCode: Succeeded, Reserved: Rsv, AddrType: Ipv4, BindAddr: make([]byte, 4), BindPort: 0}
	size = copy(bytes, ssr.ToBytes())
	bytes = bytes[:size]

	_, err = s.lstream.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) authenticate() error {
	return nil
}

func (s *Session) pipe() error {
	lch := make(chan error)
	rch := make(chan error)

	defer func() {
		close(lch)
		close(rch)
	}()

	lbuf := make([]byte, 1000)
	rbuf := make([]byte, 1000)

	pipe := func(buf []byte, readStream, writeStream net.Conn, waitCannel chan error) {
		for {
			_, err := readStream.Read(buf)
			if err != nil {
				waitCannel <- err
			}
			_, err = writeStream.Write(buf)
			if err != nil {
				waitCannel <- err
			}
		}
	}

	go pipe(lbuf, s.lstream, s.rstream, lch)
	go pipe(rbuf, s.rstream, s.lstream, rch)

	lerr := <-lch
	rerr := <-rch

	if lerr != nil {
		return lerr
	}
	if rerr != nil {
		return rerr
	}

	return nil
}

func (s *Session) shutdown() error {
	return nil
}
