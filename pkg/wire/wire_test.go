package tftp

import (
	"reflect"
	"testing"
)

func TestSerializationDeserialization(t *testing.T) {
	tests := []struct {
		bytes  []byte
		packet Packet
	}{
		{
			[]byte("\x00\x01foo\x00bar\x00"),
			&PacketRequest{OpRRQ, "foo", "bar"},
		},
		{
			[]byte("\x00\x02foo\x00bar\x00"),
			&PacketRequest{OpWRQ, "foo", "bar"},
		},
		{
			[]byte("\x00\x03\x12\x34fnord"),
			&PacketData{0x1234, []byte("fnord")},
		},
		{
			[]byte("\x00\x03\x12\x34"),
			&PacketData{0x1234, []byte("")},
		},
		{
			[]byte("\x00\x04\xd0\x0f"),
			&PacketAck{0xd00f},
		},
		{
			[]byte("\x00\x05\xab\xcdparachute failure\x00"),
			&PacketError{0xabcd, "parachute failure"},
		},
	}

	for _, test := range tests {
		actualBytes := test.packet.Serialize()
		if !reflect.DeepEqual(test.bytes, actualBytes) {
			t.Errorf("Serializing %#v: expected %q; got %q", test.packet, test.bytes, actualBytes)
		}

		actualPacket, err := ParsePacket(test.bytes)
		if err != nil {
			t.Errorf("Unable to parse packet %q: %s", test.bytes, err)
		} else if !reflect.DeepEqual(test.packet, actualPacket) {
			t.Errorf("Deserializing %q: expected %#v; got %#v", test.bytes, test.packet, actualPacket)
		}
	}
}

func TestDeserializationInvalid(t *testing.T) {
	tests := [][]byte{
		// no opcode
		[]byte(""),

		// invalid opcode
		[]byte("\x00\x00"),
		[]byte("\x00\x06"),
		[]byte("\xff\x01"),
		[]byte("\xff\xff"),

		// short RRQ
		[]byte("\x00\x01"),
		[]byte("\x00\x01foo"),
		[]byte("\x00\x01foo\x00"),
		[]byte("\x00\x01foo\x00bar"),

		// short WRQ
		[]byte("\x00\x02"),
		[]byte("\x00\x02foo"),
		[]byte("\x00\x02foo\x00"),
		[]byte("\x00\x02foo\x00bar"),

		// short data
		[]byte("\x00\x03"),
		[]byte("\x00\x03\x01"),

		// short ack
		[]byte("\x00\x04"),
		[]byte("\x00\x04\x01"),

		// short error
		[]byte("\x00\x05"),
		[]byte("\x00\x05\xab"),
		[]byte("\x00\x05\xab\xcd"),
		[]byte("\x00\x05\xab\xcdparachute failure"),
	}

	for _, test := range tests {
		if p, err := ParsePacket(test); err == nil {
			t.Errorf("Parsing packet %q: expected error; got %#v", test, p)
		}
	}
}
