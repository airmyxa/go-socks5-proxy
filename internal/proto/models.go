package proto

import "fmt"

type ParseError struct {
	message string
}

func (pe ParseError) Error() string {
	return pe.message
}

type FromBytesConstructer interface {
	FromBytes(bytes []byte) error
}

// SuggestedMethods

type SuggestedMethods struct {
	Ver      byte
	NMethods byte
	Methods  [255]byte
}

func (sm *SuggestedMethods) FromBytes(bytes []byte) error {
	if len(bytes) < 3 {
		return ParseError{message: fmt.Sprintf("invalid length of message: %d", len(bytes))}
	}

	sm.Ver = bytes[0]
	if sm.Ver != 0x5 {
		return ParseError{message: fmt.Sprintf("invalid protocol number: %d", sm.Ver)}
	}

	sm.NMethods = bytes[1]
	if sm.NMethods == 0 {
		return ParseError{fmt.Sprintf("invalid number of methods: %d", sm.NMethods)}
	}

	copy(sm.Methods[:], bytes[2:2+sm.NMethods])

	return nil
}

//SelectedMethods

type SelectedMethods struct {
	Ver    byte
	Method byte
}
