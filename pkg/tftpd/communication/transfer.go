package communication

import (
	"fmt"
	"net"
	"time"

	"github.com/BenWhiting/go-tftp/internal/constants"
	"github.com/BenWhiting/go-tftp/pkg/wire"
)

const maxDataSize = 512

// Transfer struct for active transfer
type Transfer struct {
	LastOp   time.Time
	Addr     net.Addr
	Conn     net.PacketConn
	BlockNum uint16
	Filename string
	Data     []byte
	SendData []byte
	Retry    bool
}

// Send and Transmit a serialized transfer
func (t *Transfer) Send(p wire.Packet) error {
	t.SendData = p.Serialize()
	return t.Transmit()
}

// Transmit the transfer
func (t *Transfer) Transmit() error {
	_, err := t.Conn.WriteTo(t.SendData, t.Addr)
	if nil != err {
		t.Retry = true
		return fmt.Errorf(constants.FailedSendPacketMsg, err)
	}
	t.Retry = false
	t.LastOp = time.Now()
	return nil
}
