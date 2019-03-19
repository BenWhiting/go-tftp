package communication

import (
	"testing"

	"github.com/BenWhiting/go-tftp/pkg/wire"
)

func TestTransfer_Send(t *testing.T) {
	type args struct {
		p wire.Packet
	}
	tests := []struct {
		name    string
		t       *Transfer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Send(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("Transfer.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransfer_Transmit(t *testing.T) {
	tests := []struct {
		name    string
		t       *Transfer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Transmit(); (err != nil) != tt.wantErr {
				t.Errorf("Transfer.Transmit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
