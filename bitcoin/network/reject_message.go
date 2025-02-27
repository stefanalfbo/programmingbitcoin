package network

import "io"

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
	reason string
	// Optional extra data provided by some errors. Currently, all errors which
	// provide this field fill it with the TXID or block header hash of the
	// object being rejected, so the field is 32 bytes.
	data []byte
}

// NewRejectMessage returns a new RejectMessage
func NewRejectMessage(message string, ccode CCode, reason string, data []byte) *RejectMessage {
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
	// Serialize the message
	messageBytes := []byte(rm.Message)
	messageLen := len(messageBytes)
	if messageLen > 12 {
		messageLen = 12
	}
	messageBytes = messageBytes[:messageLen]

	// Serialize the reason
	reasonBytes := []byte(rm.reason)
	reasonLen := len(reasonBytes)
	if reasonLen > 12 {
		reasonLen = 12
	}
	reasonBytes = reasonBytes[:reasonLen]

	// Serialize the data
	dataLen := len(rm.data)
	if dataLen > 32 {
		dataLen = 32
	}
	data := rm.data[:dataLen]

	// Serialize the reject message
	serialized := make([]byte, 1+1+1+len(messageBytes)+1+len(reasonBytes)+1+32)
	serialized[0] = byte(len(messageBytes))
	serialized[1] = byte(rm.CCode)
	serialized[2] = byte(len(reasonBytes))
	copy(serialized[3:], messageBytes)
	copy(serialized[3+len(messageBytes)+1:], reasonBytes)
	copy(serialized[3+len(messageBytes)+1+len(reasonBytes)+1:], data)

	return serialized, nil
}

// Parse deserializes the RejectMessage from a reader
func (rm *RejectMessage) Parse(reader io.Reader) (Message, error) {
	messageLenBytes := make([]byte, 1)
	_, err := reader.Read(messageLenBytes)
	if err != nil {
		return nil, err
	}
	messageLen := int(messageLenBytes[0])

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
	reasonLenBytes := make([]byte, 1)
	_, err = reader.Read(reasonLenBytes)
	if err != nil {
		return nil, err
	}
	reasonLen := int(reasonLenBytes[0])

	reasonBytes := make([]byte, reasonLen)
	_, err = reader.Read(reasonBytes)
	if err != nil {
		return nil, err
	}
	rm.reason = string(reasonBytes)

	data := make([]byte, 32)
	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}
	rm.data = data

	return rm, nil
}
