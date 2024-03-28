package proto

// https://www.rfc-editor.org/rfc/rfc1928

import (
	"fmt"
	"slices"
)

type Cmd byte

const (
	Connect Cmd = 0x1
	Bind    Cmd = 0x2
	Udp     Cmd = 0x3
)

var Commands = [...]Cmd{Connect, Bind, Udp}

const Ver5 = 0x5
const Rsv = 0x0

type Atyp byte

const (
	Ipv4   Atyp = 0x1
	Domain Atyp = 0x3
	Ipv6   Atyp = 0x6
)

var Atyps = [...]Atyp{Ipv4, Domain, Ipv6}

type Rep byte

const (
	Succeeded               Rep = 0x0
	GeneralServerFailure    Rep = 0x1
	ConnectionNotAllowed    Rep = 0x2
	NetworkUnreachable      Rep = 0x3
	HostUnreachable         Rep = 0x4
	ConnectionRefused       Rep = 0x5
	TtlExpired              Rep = 0x6
	CommandNotSupported     Rep = 0x7
	AddressTypeNotSupported Rep = 0x8
)

type ParseError struct {
	message string
}

func (pe ParseError) Error() string {
	return pe.message
}

type FromBytesReader interface {
	FromBytes(bytes []byte) error
}

type ToBytesWriter interface {
	ToBytes() []byte
}

// ClientSuggestedMethods

type ClientSuggestedMethods struct {
	Ver      byte
	NMethods byte
	Methods  [255]byte
}

func (sm *ClientSuggestedMethods) FromBytes(bytes []byte) error {
	if len(bytes) < 3 {
		return ParseError{message: fmt.Sprintf("invalid length of message: %d", len(bytes))}
	}

	sm.Ver = bytes[0]
	if sm.Ver != Ver5 {
		return ParseError{message: fmt.Sprintf("invalid protocol version: %d", sm.Ver)}
	}

	sm.NMethods = bytes[1]
	if sm.NMethods == 0 {
		return ParseError{fmt.Sprintf("invalid number of methods: %d", sm.NMethods)}
	}

	copy(sm.Methods[:], bytes[2:2+sm.NMethods])

	return nil
}

//ServerSelectedMethod

type ServerSelectedMethod struct {
	Ver    byte
	Method byte
}

func (sm *ServerSelectedMethod) ToBytes() []byte {
	bytes := make([]byte, 2)
	bytes[0] = sm.Ver
	bytes[1] = sm.Method
	return bytes
}

// ClientSocksRequest

type ClientSocksRequest struct {
	Ver      byte
	Command  Cmd
	Reserved byte
	AddrType Atyp
	DestAddr []byte
	DestPort uint16
}

func (sr *ClientSocksRequest) FromBytes(bytes []byte) error {
	if len(bytes) < 6 {
		return ParseError{message: fmt.Sprintf("Incorrent length of client socks request: %d", len(bytes))}
	}

	sr.Ver = bytes[0]
	if sr.Ver != Ver5 {
		return ParseError{message: fmt.Sprintf("Incorrent version of protocol: %b", sr.Ver)}
	}

	sr.Command = Cmd(bytes[1])
	if !slices.Contains(Commands[:], sr.Command) {
		return ParseError{message: fmt.Sprintf("Unsupported command: %b", sr.Command)}
	}

	sr.Reserved = bytes[2]
	if sr.Reserved != Rsv {
		return ParseError{message: fmt.Sprintf("Incorrect reserve field: %b", sr.Reserved)}
	}

	sr.AddrType = Atyp(bytes[3])
	if !slices.Contains(Atyps[:], sr.AddrType) {
		return ParseError{message: fmt.Sprintf("Unsupported address type: %b", sr.AddrType)}
	}

	var last uint8
	if sr.AddrType == Ipv4 {
		sr.DestAddr = make([]byte, 4)
		copy(sr.DestAddr, bytes[4:8])
		last = 8
	} else if sr.AddrType == Ipv6 {
		sr.DestAddr = make([]byte, 20)
		copy(sr.DestAddr, bytes[4:20])
		last = 20
	} else if sr.AddrType == Domain {
		length := bytes[4]
		sr.DestAddr = make([]byte, length)
		copy(sr.DestAddr, bytes[5:5+length])
		last = length
	} else {
		return ParseError{message: fmt.Sprintf("Cannot parse destination address of type: %b", sr.AddrType)}
	}

	sr.DestPort = uint16(bytes[last])
	sr.DestPort <<= 8
	sr.DestPort |= uint16(bytes[last+1])

	return nil
}

type ServerSocksResponse struct {
	Ver       byte
	ReplyCode Rep
	Reserved  byte
	AddrType  Atyp
	BindAddr  []byte
	BindPort  uint16
}

func (sr *ServerSocksResponse) ToBytes() []byte {
	bytes := make([]byte, 10)

	bytes[0] = sr.Ver
	bytes[1] = byte(sr.ReplyCode)
	bytes[2] = sr.Reserved
	bytes[3] = byte(sr.AddrType)

	copy(bytes[4:], sr.BindAddr)
	last := 4 + len(sr.BindAddr)

	bytes[last] = byte(sr.BindPort >> 8)
	bytes[last+1] = byte(sr.BindPort & 255)

	return bytes
}
