package tftpd

import (
	"reflect"
	"testing"
)

func TestNewTFTPServer(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTFTPServer(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTFTPServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTFTPServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
