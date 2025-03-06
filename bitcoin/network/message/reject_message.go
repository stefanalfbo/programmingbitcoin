package message

import (
	"io"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type CCode byte

const (
	// Protocol syntax error
	REJECT_MALFORMED = 0x01
	// Protocol semantic error
	REJECT_INVALID   = 0x10
	REJECT_OBSOLETE  = 0x11
	REJECT_DUPLICATE = 0x12
	// Ser policy rule
	REJECT_NONSTANDARD     = 0x40
	REJECT_DUST            = 0x41
	REJECT_INSUFFICIENTFEE = 0x42
	REJECT_CHECKPOINT      = 0x43
)

// The reject message is sent when messages are rejected.
type RejectMessage struct {
	// type of message rejected
	Message string
	// code relating to rejected message
	CCode CCode
	// text version of reason for rejection
	Reason string
	// Optional extra data provided by some errors. Currently, all errors which
	// provide this field fill it with the TXID or block header hash of the
	// object being rejected, so the field is 32 bytes.
	Data [32]byte
}

// NewRejectMessage returns a new RejectMessage
func NewRejectMessage(message string, ccode CCode, reason string, data [32]byte) *RejectMessage {
	return &RejectMessage{message, ccode, reason, data}
}

// NewEmptyRejectMessage returns a new RejectMessage with empty fields
func NewEmptyRejectMessage() *RejectMessage {
	return &RejectMessage{}
}

// Command returns the protocol command string for the message
func (rm *RejectMessage) Command() []byte {
	return []byte("reject")
}

// Serialize serializes the RejectMessage
func (rm *RejectMessage) Serialize() ([]byte, error) {
	serialized := make([]byte, 0)

	messageLen, err := varint.Encode(uint64(len(rm.Message)))
	if err != nil {
		return nil, err
	}
	serialized = append(serialized, messageLen...)
	serialized = append(serialized, []byte(rm.Message)...)

	serialized = append(serialized, byte(rm.CCode))

	reasonLen, err := varint.Encode(uint64(len(rm.Reason)))
	if err != nil {
		return nil, err
	}
	serialized = append(serialized, reasonLen...)
	serialized = append(serialized, []byte(rm.Reason)...)

	serialized = append(serialized, rm.Data[:]...)

	return serialized, nil
}

// Parse deserializes the RejectMessage from a reader
func (rm *RejectMessage) Parse(reader io.Reader) (Message, error) {
	messageLen, err := varint.Decode(reader)
	if err != nil {
		return nil, err
	}

	messageBytes := make([]byte, messageLen)
	_, err = reader.Read(messageBytes)
	if err != nil {
		return nil, err
	}
	rm.Message = string(messageBytes)

	// Read the ccode
	ccodeBytes := make([]byte, 1)
	_, err = reader.Read(ccodeBytes)
	if err != nil {
		return nil, err
	}
	rm.CCode = CCode(ccodeBytes[0])

	// Read the reason
	reasonLen, err := varint.Decode(reader)
	if err != nil {
		return nil, err
	}

	reasonBytes := make([]byte, reasonLen)
	_, err = reader.Read(reasonBytes)
	if err != nil {
		return nil, err
	}
	rm.Reason = string(reasonBytes)

	data := make([]byte, 32)
	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}
	copy(rm.Data[:], data[:32])

	return rm, nil
}
