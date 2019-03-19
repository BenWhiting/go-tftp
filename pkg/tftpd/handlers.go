package tftpd

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/BenWhiting/go-tftp/internal/constants"
	"github.com/BenWhiting/go-tftp/pkg/tftpd/communication"
	"github.com/BenWhiting/go-tftp/pkg/wire"
)

/*
  2 bytes    string   1 byte     string   1 byte
 -----------------------------------------------
 |  01  |  Filename  |   0  |    Mode    |   0  |
 -----------------------------------------------
 Mode = octet
 NOTE: Will only handle octet mode
*/
func (s *Server) rrqHandler(conn net.PacketConn, p wire.Packet, addr net.Addr) {
	// Cast packet
	pkt := p.(*wire.PacketRequest)

	// Log request occurrence
	s.loggerToFile.Printf(constants.ReceivedMsg, pkt.Op, p)

	// Check Mode type
	if s.Config.AllowedMode == strings.ToLower(pkt.Mode) {
		msg := fmt.Sprintf(constants.ModeNotSupportedMsg, pkt.Op, pkt.Mode)
		serverMsg := fmt.Sprintf(constants.ServerMessageWrapper, msg, addr)
		s.logger.Println(serverMsg)
		s.loggerToFile.Println(serverMsg)
		s.sendError(0, addr, conn, msg)
		return
	}

	// Check file exists
	s.mux.Lock()
	fileData, ok := s.store[pkt.Filename]
	s.mux.Unlock()
	if !ok {
		msg := fmt.Sprintf(constants.RRQFileNotFoundErrMsg, pkt.Filename)
		serverMsg := fmt.Sprintf(constants.ServerMessageWrapper, msg, addr)
		s.logger.Println(serverMsg)
		s.loggerToFile.Println(serverMsg)
		s.sendError(1, addr, conn, msg)
		return
	}

	// Create Transfer
	t := &communication.Transfer{
		BlockNum: 1,
		Addr:     addr,
		Conn:     conn,
		Data:     fileData,
		Filename: pkt.Filename,
	}

	// Break down message based on MaxDataSize
	var data []byte
	if len(t.Data) >= constants.MaxDataSize {
		data = t.Data[:constants.MaxDataSize]
	} else {
		data = t.Data
	}

	// Create response package
	res := wire.PacketData{
		Data:     data,
		BlockNum: t.BlockNum}

	// Send
	t.Send(&res)

	// Add file to active transfers
	s.mux.Lock()
	s.activeTransfers[t.Filename] = t
	s.mux.Unlock()
}

/*
  2 bytes    string   1 byte     string   1 byte
 -----------------------------------------------
 |  02 |  Filename  |   0  |    Mode    |   0  |
 -----------------------------------------------
 Mode = octet
 NOTE: Will only handle octet  mode
*/

func (s *Server) wrqHandler(conn net.PacketConn, p wire.Packet, addr net.Addr) {
	// Cast packet
	pkt := p.(*wire.PacketRequest)

	// Log request occurrence
	s.loggerToFile.Printf(constants.ReceivedMsg, pkt.Op, p)

	// Check Mode type
	if s.Config.AllowedMode == strings.ToLower(pkt.Mode) {
		msg := fmt.Sprintf(constants.ModeNotSupportedMsg, pkt.Op, pkt.Mode)
		serverMsg := fmt.Sprintf(constants.ServerMessageWrapper, msg, addr)
		s.logger.Println(serverMsg)
		s.loggerToFile.Println(serverMsg)
		s.sendError(0, addr, conn, msg)
		return
	}

	// Check file exists
	s.mux.Lock()
	_, ok := s.store[pkt.Filename]
	s.mux.Unlock()
	if !ok {
		msg := fmt.Sprintf(constants.WRQFileAlreadyExistsErrMsg, pkt.Filename)
		serverMsg := fmt.Sprintf(constants.ServerMessageWrapper, msg, addr)
		s.logger.Println(serverMsg)
		s.loggerToFile.Println(serverMsg)
		s.sendError(6, addr, conn, msg)
		return
	}

	// Create Transfer
	t := &communication.Transfer{
		BlockNum: 1,
		Addr:     addr,
		Conn:     conn,
		Filename: pkt.Filename,
	}

	// Create response package
	res := wire.PacketAck{BlockNum: 0}

	// Send
	t.Send(&res)

	// Add file to active transfers
	s.mux.Lock()
	s.activeTransfers[t.Filename] = t
	s.mux.Unlock()
}

/*
 2 bytes    2 bytes       n bytes
 ---------------------------------
| 03    |   Block #  |    Data    |
 ---------------------------------
*/
func (s *Server) dataHandler(conn net.PacketConn, p wire.Packet, addr net.Addr, n int) {
	// Cast packet
	pkt := p.(*wire.PacketData)

	// Find active transfer
	t, err := s.findTransfer(pkt.BlockNum, addr)
	if err != nil {
		s.logger.Printf(constants.SimpleErrMsg, err)
		s.sendError(5, addr, conn, "")
		return
	}

	// Trim NULL characters
	d := bytes.Trim(pkt.Data, "\x00")

	// The buffer is reused, so leftovers have to be cut out
	t.Data = append(t.Data, d[:n]...)
	t.BlockNum++

	// Create response package and send
	res := wire.PacketAck{BlockNum: pkt.BlockNum}
	t.Send(&res)

	// If incoming size is less than max data size file is done
	// else continue
	s.mux.Lock()
	if n < constants.MaxDataSize {
		s.store[t.Filename] = t.Data
		delete(s.activeTransfers, t.Filename)
		s.logger.Printf(constants.FileReceivingCompleteMsg, t.Filename, addr)
	} else {
		s.activeTransfers[t.Filename] = t
	}
	s.mux.Unlock()

}

/*
  2 bytes     2 bytes
 ---------------------
|   04  |   Block #  |
 ---------------------
*/
func (s *Server) ackHandler(conn net.PacketConn, p wire.Packet, addr net.Addr) {
	// Cast packet
	pkt := p.(*wire.PacketAck)

	// Find the file transfer
	t, err := s.findTransfer(pkt.BlockNum, addr)
	if err != nil {
		return
	}

	// Move to the next block
	t.BlockNum++

	// Find the transfer size and read size
	tSize := t.BlockNum * constants.MaxDataSize
	rSize := pkt.BlockNum * constants.MaxDataSize

	var data []byte
	if len(t.Data) >= int(tSize) {
		data = t.Data[rSize:tSize]
	} else if len(t.Data) < int(rSize) {
		// Leaving the transfer in activeTransfers in case of pending retransmits.
		s.logger.Printf(constants.FileTransferCompleteMsg, t.Filename)
		s.mux.Lock()
		delete(s.activeTransfers, t.Filename)
		s.mux.Unlock()
		return
	} else {
		data = t.Data[rSize:]
	}

	// Create response package and send
	res := wire.PacketData{BlockNum: t.BlockNum, Data: data}
	t.Send(&res)

	s.mux.Lock()
	s.activeTransfers[t.Filename] = t
	s.mux.Unlock()
}

/*
 2 bytes     2 bytes      string    1 byte
 -----------------------------------------
|  05   |  ErrorCode |   ErrMsg   |   0  |
 -----------------------------------------
*/
func (s *Server) errHandler(p wire.Packet, addr net.Addr) {
	pkt := p.(*wire.PacketError)
	s.logger.Printf(constants.ReceivedErrMsg, addr, pkt.Code, pkt.Msg)
}

func (s *Server) findTransfer(block uint16, addr net.Addr) (t *communication.Transfer, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, t := range s.activeTransfers {
		if t.BlockNum == block && addr.String() == t.Addr.String() {
			return t, nil
		}
	}

	return nil, errors.New(constants.UnknownTransferMsg)
}
