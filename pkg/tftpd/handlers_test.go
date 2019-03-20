package tftpd

import (
	"net"
	"os"
	"testing"

	"github.com/BenWhiting/go-tftp/pkg/wire"
)

func TestServer_rrqHandler(t *testing.T) {

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
	defer s.Connection.Close()

	type args struct {
		conn net.PacketConn
		p    wire.Packet
		addr net.Addr
	}
	tests := []struct {
		name        string
		server      *Server
		args        args
		expectedErr string
		wantErr     bool
	}{
		{name: "File not found",
			server:      s,
			args:        args{addr: cli.LocalAddr(), conn: s.Connection, p: &wire.PacketRequest{Op: 1, Filename: "testfile.txt", Mode: "octet"}},
			expectedErr: "[RRQ] File testfile.txt was not found. ",
			wantErr:     true},
		{name: "Wrong mode found",
			server:      s,
			args:        args{addr: cli.LocalAddr(), conn: s.Connection, p: &wire.PacketRequest{Op: 1, Filename: "testfile.txt", Mode: "FAKE MODE"}},
			expectedErr: "[1] Mode FAKE MODE is not supported.",
			wantErr:     true},
		{name: "Valid Input",
			server:      s,
			args:        args{addr: cli.LocalAddr(), conn: s.Connection, p: &wire.PacketRequest{Op: 1, Filename: "testfile.txt", Mode: "octet"}},
			expectedErr: "",
			wantErr:     false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.rrqHandler(tt.args.conn, tt.args.p, tt.args.addr)

			out := make([]byte, 1024)
			_, err := cli.Read(out)
			if nil != err {
				t.Errorf(err.Error())
			}
			p, err := wire.ParsePacket(out)
			if nil != err {
				t.Errorf(err.Error())
			}
			if tt.wantErr {
				pkt := p.(*wire.PacketError)
				if pkt.Msg != tt.expectedErr {
					t.Errorf("Unexpected return error")
				}
			}
		})
		s.logFile.Close()
		os.Remove(s.Config.LogFilePath)
	}
}

func TestServer_wrqHandler(t *testing.T) {
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
	defer s.Connection.Close()

	type args struct {
		conn net.PacketConn
		p    wire.Packet
		addr net.Addr
	}
	tests := []struct {
		name        string
		server      *Server
		args        args
		expectedErr string
		wantErr     bool
	}{
		{name: "Write file to server",
			server: s,
			args:   args{addr: cli.LocalAddr(), conn: s.Connection, p: &wire.PacketRequest{Op: 2, Filename: "testfile.txt", Mode: "octet"}}},
		{name: "Wrong mode found",
			server:      s,
			args:        args{addr: cli.LocalAddr(), conn: s.Connection, p: &wire.PacketRequest{Op: 2, Filename: "testfile.txt", Mode: "BADMODE"}},
			expectedErr: "[2] Mode BADMODE is not supported.",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.wrqHandler(tt.args.conn, tt.args.p, tt.args.addr)
			out := make([]byte, 1024)
			_, err := cli.Read(out)
			p, err := wire.ParsePacket(out)
			if nil != err {
				t.Errorf(err.Error())
			}
			if tt.wantErr {
				pkt := p.(*wire.PacketError)
				if pkt.Msg != tt.expectedErr {
					t.Errorf("Unexpected return error")
				}
			}
		})
		s.logFile.Close()
		os.Remove(s.Config.LogFilePath)
	}
}
