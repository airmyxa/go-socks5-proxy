package proto

import "net"

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
	return nil
}

func (s *Session) handshake() error {
	return nil
}

func (s *Session) authenticate() error {
	return nil
}

func (s *Session) pipe() error {
	return nil
}

func (s *Session) shutdown() error {
	return nil
}
