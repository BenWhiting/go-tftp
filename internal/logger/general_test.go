package logger

import (
	"testing"
)

func TestNewGeneralLogger(t *testing.T) {
	type args struct {
		flag int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Check under 1 is not allowed", args: args{flag: 0}, wantErr: true},
		{name: "Check above 7 is not allowed", args: args{flag: 8}, wantErr: true},
		{name: "Check 1 is allowed", args: args{flag: 1}, wantErr: false},
		{name: "Check 7 is allowed", args: args{flag: 7}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGeneralLogger(tt.args.flag)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGeneralLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
