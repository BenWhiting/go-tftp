package tftpd

import (
	"net"
	"reflect"
	"testing"

	"github.com/BenWhiting/go-tftp/pkg/tftpd/communication"
	"github.com/BenWhiting/go-tftp/pkg/wire"
)

func TestServer_rrqHandler(t *testing.T) {
	type args struct {
		conn net.PacketConn
		p    wire.Packet
		addr net.Addr
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.rrqHandler(tt.args.conn, tt.args.p, tt.args.addr)
		})
	}
}

func TestServer_wrqHandler(t *testing.T) {
	type args struct {
		conn net.PacketConn
		p    wire.Packet
		addr net.Addr
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.wrqHandler(tt.args.conn, tt.args.p, tt.args.addr)
		})
	}
}

func TestServer_dataHandler(t *testing.T) {
	type args struct {
		conn net.PacketConn
		p    wire.Packet
		addr net.Addr
		n    int
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.dataHandler(tt.args.conn, tt.args.p, tt.args.addr, tt.args.n)
		})
	}
}

func TestServer_ackHandler(t *testing.T) {
	type args struct {
		conn net.PacketConn
		p    wire.Packet
		addr net.Addr
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.ackHandler(tt.args.conn, tt.args.p, tt.args.addr)
		})
	}
}

func TestServer_errHandler(t *testing.T) {
	type args struct {
		p    wire.Packet
		addr net.Addr
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.errHandler(tt.args.p, tt.args.addr)
		})
	}
}

func TestServer_findTransfer(t *testing.T) {
	type args struct {
		block uint16
		addr  net.Addr
	}
	tests := []struct {
		name    string
		s       *Server
		args    args
		wantT   *communication.Transfer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := tt.s.findTransfer(tt.args.block, tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.findTransfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("Server.findTransfer() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}
