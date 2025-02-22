package network

import (
	"encoding/binary"
	"io"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

// When a node creates an outgoing connection, it will immediately
// advertise its version. The remote node will respond with its
// version. No further communication is possible until both peers
// have exchanged their version. See:
// https://en.bitcoin.it/wiki/Protocol_documentation#version
type VersionMessage struct {
	command []byte
	// Identifies protocol version being used by the node
	Version int32
	// Bit field of features to be enabled for this connection
	Services uint64
	// Standard UNIX timestamp in seconds
	Timestamp int64
	// The network address of the node receiving this message
	ReceiverServices uint64
	ReceiverIP       [16]byte
	ReceiverPort     uint16
	// Fields below require version ≥ 106
	// Field can be ignored. This used to be the network address of
	// the node emitting this message, but most P2P implementations
	// send 26 dummy bytes. The "services" field of the address would
	// also be redundant with the second field of the version message.
	SenderServices uint64
	SenderIP       [16]byte
	SenderPort     uint16
	// Node random nonce, randomly generated every time a version packet
	// is sent. This nonce is used to detect connections to self.
	Nonce uint64
	// User Agent (0x00 if string is 0 bytes long)
	UserAgent []byte
	// The last block received by the emitting node
	LatestBlock int32
	// Fields below require version ≥ 70001
	// Whether the remote peer should announce relayed transactions or not, see:
	// BIP 0037 - https://github.com/bitcoin/bips/blob/master/bip-0037.mediawiki
	Relay bool
}

// NewVersionMessage returns a new VersionMessage
func NewVersionMessage() *VersionMessage {
	return &VersionMessage{
		command:          []byte("version"),
		Version:          70015,
		Services:         0,
		Timestamp:        0,
		ReceiverServices: 0,
		ReceiverIP:       [16]byte{0x00, 0x00, 0x00, 0x00},
		ReceiverPort:     8333,
		SenderServices:   0,
		SenderIP:         [16]byte{0x00, 0x00, 0x00, 0x00},
		SenderPort:       8333,
		Nonce:            0,
		UserAgent:        []byte("/programmingbitcoin:0.1/"),
		LatestBlock:      0,
		Relay:            false,
	}
}

// Command returns the command of a VersionMessage
func (vm *VersionMessage) Command() []byte {
	return vm.command
}

// Serialize serializes a VersionMessage
func (vm *VersionMessage) Serialize() ([]byte, error) {
	ip4Prefix := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}
	result := make([]byte, 0)

	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(vm.Version))
	result = append(result, version...)

	services := make([]byte, 8)
	binary.LittleEndian.PutUint64(services, vm.Services)
	result = append(result, services...)

	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(vm.Timestamp))
	result = append(result, timestamp...)

	receiverServices := make([]byte, 8)
	binary.LittleEndian.PutUint64(receiverServices, vm.ReceiverServices)
	result = append(result, receiverServices...)

	result = append(result, ip4Prefix...)
	result = append(result, vm.ReceiverIP[:4]...)

	h, l := uint(vm.ReceiverPort)>>8, uint(vm.ReceiverPort)&0xff
	result = append(result, byte(h), byte(l))

	senderServices := make([]byte, 8)
	binary.LittleEndian.PutUint64(senderServices, vm.SenderServices)
	result = append(result, senderServices...)

	result = append(result, ip4Prefix...)
	result = append(result, vm.SenderIP[:4]...)

	h, l = uint(vm.SenderPort)>>8, uint(vm.SenderPort)&0xff
	result = append(result, byte(h), byte(l))

	nonce := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonce, vm.Nonce)
	result = append(result, nonce...)

	userAgentLength, err := varint.Encode(uint64(len(vm.UserAgent)))
	if err != nil {
		return nil, err
	}

	result = append(result, userAgentLength...)
	result = append(result, vm.UserAgent...)

	latestBlock := make([]byte, 4)
	binary.LittleEndian.PutUint32(latestBlock, uint32(vm.LatestBlock))
	result = append(result, latestBlock...)

	if vm.Relay {
		result = append(result, 0x01)
	} else {
		result = append(result, 0x00)
	}

	return result, nil
}

// Parse parses a VersionMessage
func (vm *VersionMessage) Parse(reader io.Reader) (Message, error) {
	message := VersionMessage{}

	var version int32
	err := binary.Read(reader, binary.LittleEndian, &version)
	if err != nil {
		return nil, err
	}
	message.Version = version

	var services uint64
	err = binary.Read(reader, binary.LittleEndian, &services)
	if err != nil {
		return nil, err
	}
	message.Services = services

	var timestamp int64
	err = binary.Read(reader, binary.LittleEndian, &timestamp)
	if err != nil {
		return nil, err
	}
	message.Timestamp = timestamp

	var receiverServices uint64
	err = binary.Read(reader, binary.LittleEndian, &receiverServices)
	if err != nil {
		return nil, err
	}
	message.ReceiverServices = receiverServices

	receiverIP := [16]byte{}
	err = binary.Read(reader, binary.LittleEndian, &receiverIP)
	if err != nil {
		return nil, err
	}
	message.ReceiverIP = receiverIP

	var receiverPort uint16
	err = binary.Read(reader, binary.LittleEndian, &receiverPort)
	if err != nil {
		return nil, err
	}
	message.ReceiverPort = receiverPort

	var senderServices uint64
	err = binary.Read(reader, binary.LittleEndian, &senderServices)
	if err != nil {
		return nil, err
	}
	message.SenderServices = senderServices

	senderIP := [16]byte{}
	err = binary.Read(reader, binary.LittleEndian, &senderIP)
	if err != nil {
		return nil, err
	}
	message.SenderIP = senderIP

	var senderPort uint16
	err = binary.Read(reader, binary.LittleEndian, &senderPort)
	if err != nil {
		return nil, err
	}
	message.SenderPort = senderPort

	var nonce uint64
	err = binary.Read(reader, binary.LittleEndian, &nonce)
	if err != nil {
		return nil, err
	}
	message.Nonce = nonce

	userAgentLength, err := varint.Decode(reader)
	if err != nil {
		return nil, err
	}

	userAgent := make([]byte, userAgentLength)
	err = binary.Read(reader, binary.LittleEndian, &userAgent)
	if err != nil {
		return nil, err
	}
	message.UserAgent = userAgent

	var latestBlock int32
	err = binary.Read(reader, binary.LittleEndian, &latestBlock)
	if err != nil {
		return nil, err
	}
	message.LatestBlock = latestBlock

	var relay byte
	err = binary.Read(reader, binary.LittleEndian, &relay)
	if err != nil {
		return nil, err
	}
	message.Relay = relay == 0x01

	return &message, nil
}
