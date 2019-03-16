package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// larger than a typical mtu (1500), and largest DATA packet (516).
// may limit the length of filenames in RRQ/WRQs -- RFC1350 doesn't offer a bound for these.
const MaxPacketSize = 2048

const (
	OpRRQ   uint16 = 1
	OpWRQ          = 2
	OpData         = 3
	OpAck          = 4
	OpError        = 5
)

// packet is the interface met by all packet structs
type Packet interface {
	// Parse parses a packet from its wire representation
	Parse([]byte) error
	// Serialize serializes a packet to its wire representation
	Serialize() []byte
}

// PacketRequest represents a request to read or rite a file.
type PacketRequest struct {
	Op       uint16 // OpRRQ or OpWRQ
	Filename string
	Mode     string
}

func (p *PacketRequest) Parse(buf []byte) (err error) {
	if p.Op, buf, err = parseUint16(buf); err != nil {
		return err
	}
	if p.Filename, buf, err = parseString(buf); err != nil {
		return err
	}
	if p.Mode, buf, err = parseString(buf); err != nil {
		return err
	}
	return nil
}

func (p *PacketRequest) Serialize() []byte {
	buf := make([]byte, 2+len(p.Filename)+1+len(p.Mode)+1)
	binary.BigEndian.PutUint16(buf, p.Op)
	copy(buf[2:], p.Filename)
	copy(buf[2+len(p.Filename)+1:], p.Mode)
	return buf
}

// PacketData carries a block of data in a file transmission.
type PacketData struct {
	BlockNum uint16
	Data     []byte
}

func (p *PacketData) Parse(buf []byte) (err error) {
	buf = buf[2:] // skip over op
	if p.BlockNum, buf, err = parseUint16(buf); err != nil {
		return err
	}
	p.Data = buf
	return nil
}

func (p *PacketData) Serialize() []byte {
	buf := make([]byte, 4+len(p.Data))
	binary.BigEndian.PutUint16(buf, OpData)
	binary.BigEndian.PutUint16(buf[2:], p.BlockNum)
	copy(buf[4:], p.Data)
	return buf
}

// PacketAck acknowledges receipt of a data packet
type PacketAck struct {
	BlockNum uint16
}

func (p *PacketAck) Parse(buf []byte) (err error) {
	buf = buf[2:] // skip over op
	if p.BlockNum, buf, err = parseUint16(buf); err != nil {
		return err
	}
	return nil
}

func (p *PacketAck) Serialize() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf, OpAck)
	binary.BigEndian.PutUint16(buf[2:], p.BlockNum)
	return buf
}

// PacketError is sent by a peer who has encountered an error condition
type PacketError struct {
	Code uint16
	Msg  string
}

func (p *PacketError) Parse(buf []byte) (err error) {
	buf = buf[2:] // skip over op
	if p.Code, buf, err = parseUint16(buf); err != nil {
		return err
	}
	if p.Msg, buf, err = parseString(buf); err != nil {
		return err
	}
	return nil
}

func (p *PacketError) Serialize() []byte {
	buf := make([]byte, 4+len(p.Msg)+1)
	binary.BigEndian.PutUint16(buf, OpError)
	binary.BigEndian.PutUint16(buf[2:], p.Code)
	copy(buf[4:], p.Msg)
	return buf
}

// parseUint16 reads a big-endian uint16 from the beginning of buf,
// returning it along with a slice pointing at the next position in the buffer.
func parseUint16(buf []byte) (uint16, []byte, error) {
	if len(buf) < 2 {
		return 0, nil, errors.New("packet truncated")
	}
	return binary.BigEndian.Uint16(buf), buf[2:], nil
}

// parseString reads a null-terminated ASCII string from buf,
// returning it along with a slice pointing at the next position in the buffer.
func parseString(buf []byte) (string, []byte, error) {
	i := bytes.IndexByte(buf, 0)
	if i < 0 {
		return "", nil, errors.New("packet truncated")
	}
	return string(buf[:i]), buf[i+1:], nil
}

// ParsePacket parses a packet from its wire representation.
func ParsePacket(buf []byte) (p Packet, err error) {
	var opcode uint16
	if opcode, _, err = parseUint16(buf); err != nil {
		return
	}
	switch opcode {
	case OpRRQ, OpWRQ:
		p = &PacketRequest{}
	case OpData:
		p = &PacketData{}
	case OpAck:
		p = &PacketAck{}
	case OpError:
		p = &PacketError{}
	default:
		err = fmt.Errorf("unexpected opcode %d", opcode)
		return
	}
	err = p.Parse(buf)
	return
}
