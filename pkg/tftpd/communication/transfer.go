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

func (t *Transfer) Send(r wire.Packet) error {
	t.SendData = r.Serialize()
	return t.transmit()
}

func (t *Transfer) transmit() error {
	_, err := t.Conn.WriteTo(t.SendData, t.Addr)
	if nil != err {
		t.Retry = true
		return fmt.Errorf(constants.FailedSendPacketMsg, err)
	}
	t.Retry = false
	t.LastOp = time.Now()
	return nil
}
