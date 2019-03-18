package tftpd

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/BenWhiting/go-tftp/internal/constants"
	"github.com/BenWhiting/go-tftp/pkg/wire"
)

func (s *Server) route(conn net.PacketConn, addr net.Addr, buf []byte, n int) {
	p, err := wire.ParsePacket(buf)
	if nil != err {
		s.sendError(0, addr, conn, fmt.Sprintf(constants.ParsingPackErrMsg, err))
		return
	}

	op := binary.BigEndian.Uint16(buf)

	switch op {
	case wire.OpRRQ:
		s.rrqHandler(conn, p, addr)
	case wire.OpWRQ:
		s.wrqHandler(conn, p, addr)
	case wire.OpData:
		s.dataHandler(conn, p, addr, n-4)
	case wire.OpAck:
		s.ackHandler(conn, p, addr)
	case wire.OpError:
		s.errHandler(p, addr)
	default:
		msg := fmt.Sprintf(constants.UnrecOpCodeMsg, op)
		s.loggerToFile.Println(msg)
		s.sendError(0, addr, conn, msg)
	}
}

func (s *Server) sendError(code uint16, addr net.Addr, conn net.PacketConn, errMsg string) {
	s.logger.Printf(constants.SendErrMsg, errMsg, addr)

	res := wire.PacketError{Code: code, Msg: errMsg}

	_, err := conn.WriteTo(res.Serialize(), addr)
	if nil != err {
		s.logger.Printf(constants.FailedSendErrMsg, addr, err)
	}
}
