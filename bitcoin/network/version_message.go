package network

import (
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type VersionMessage struct {
	command          []byte
	Version          int32
	Services         uint64
	Timestamp        int64
	ReceiverServices uint64
	ReceiverIP       []byte
	ReceiverPort     uint16
	SenderServices   uint64
	SenderIP         []byte
	SenderPort       uint16
	Nonce            [8]byte
	UserAgent        []byte
	LatestBlock      int32
	Relay            bool
}

func NewVersionMessage() *VersionMessage {
	return &VersionMessage{
		command:          []byte("version"),
		Version:          70015,
		Services:         0,
		Timestamp:        0,
		ReceiverServices: 0,
		ReceiverIP:       []byte{0x00, 0x00, 0x00, 0x00},
		ReceiverPort:     8333,
		SenderServices:   0,
		SenderIP:         []byte{0x00, 0x00, 0x00, 0x00},
		SenderPort:       8333,
		Nonce:            [8]byte{},
		UserAgent:        []byte("/programmingbitcoin:0.1/"),
		LatestBlock:      0,
		Relay:            false,
	}
}

func (vm *VersionMessage) Serialize() ([]byte, error) {
	ip4Prefix := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}
	result := make([]byte, 0)

	result = append(result, endian.Int32ToLittleEndian(vm.Version)...)
	result = append(result, endian.Uint64ToLittleEndian(vm.Services)...)
	result = append(result, endian.Int64ToLittleEndian(vm.Timestamp)...)
	result = append(result, endian.Uint64ToLittleEndian(vm.ReceiverServices)...)
	result = append(result, ip4Prefix...)
	result = append(result, vm.ReceiverIP[:]...)

	h, l := uint(vm.ReceiverPort)>>8, uint(vm.ReceiverPort)&0xff
	result = append(result, byte(h), byte(l))

	result = append(result, endian.Uint64ToLittleEndian(vm.SenderServices)...)

	result = append(result, ip4Prefix...)
	result = append(result, vm.SenderIP[:]...)

	h, l = uint(vm.SenderPort)>>8, uint(vm.SenderPort)&0xff
	result = append(result, byte(h), byte(l))

	result = append(result, vm.Nonce[:]...)

	userAgentLength, err := varint.Encode(big.NewInt(int64(len(vm.UserAgent))))
	if err != nil {
		return nil, err
	}

	result = append(result, userAgentLength...)
	result = append(result, vm.UserAgent...)

	result = append(result, endian.Int32ToLittleEndian(vm.LatestBlock)...)

	if vm.Relay {
		result = append(result, 0x01)
	} else {
		result = append(result, 0x00)
	}

	return result, nil
}
