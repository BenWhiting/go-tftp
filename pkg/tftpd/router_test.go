package tftpd

import (
	"net"
	"os"
	"testing"

	"github.com/BenWhiting/go-tftp/pkg/wire"
)

func TestServer_sendError(t *testing.T) {
	validServerConfig := &Config{
		AllowedMode:     "octet",
		Address:         "127.0.0.1:69",
		LogFilePath:     "./test.txt",
		FlushInterval:   1,
		LogFlag:         1,
		TransferTimeout: 1,
		RetryTime:       1,
	}
	s, err := NewTFTPServer(validServerConfig)
	if nil != err {
		t.Errorf("set up failure")
		return
	}
	s.Initialize()
	if nil != err {
		t.Errorf("set up failure")
		return
	}
	go s.Start()

	cli, err := net.Dial("udp", s.Config.Address)
	if err != nil {
		t.Error("could not connect to server: ", err)
		s.logFile.Close()
		os.Remove(s.Config.LogFilePath)
		return
	}
	defer cli.Close()

	type args struct {
		code   uint16
		addr   net.Addr
		conn   net.PacketConn
		errMsg string
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		{name: "Read short message from server",
			s: s,
			args: args{code: 0, errMsg: "test", addr: cli.LocalAddr(),
				conn: s.Connection}},
		{name: "Read more complex message from server",
			s: s,
			args: args{code: 0, errMsg: "SUPERDY LONGER MESSA!_(GE SEND BECAUSE LIFE Is full of things 123124123413123",
				addr: cli.LocalAddr(),
				conn: s.Connection}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.sendError(tt.args.code, tt.args.addr, tt.args.conn, tt.args.errMsg)
			out := make([]byte, 1024)
			_, err := cli.Read(out)
			if nil != err {
				t.Errorf(err.Error())
			}
			p, err := wire.ParsePacket(out)
			if nil != err {
				t.Errorf(err.Error())
			}
			pkt := p.(*wire.PacketError)
			if pkt.Msg != tt.args.errMsg {
				t.Errorf("Message lost in send")
			}
			if pkt.Code != tt.args.code {
				t.Errorf("code lost in send")
			}
		})

		tt.s.logFile.Close()
		os.Remove(tt.s.Config.LogFilePath)
	}
}
