package tftpd

import (
	"net"
	"testing"
)

func TestServer_route(t *testing.T) {
	type args struct {
		conn net.PacketConn
		addr net.Addr
		buf  []byte
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
			tt.s.route(tt.args.conn, tt.args.addr, tt.args.buf, tt.args.n)
		})
	}
}

func TestServer_sendError(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.sendError(tt.args.code, tt.args.addr, tt.args.conn, tt.args.errMsg)
		})
	}
}
