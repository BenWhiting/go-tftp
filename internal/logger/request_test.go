package logger

import (
	"os"
	"testing"
)

func TestNewRequestLogger(t *testing.T) {
	type args struct {
		filePath string
		flag     int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Check under 1 is not allowed", args: args{flag: 0, filePath: "./test"}, wantErr: true},
		{name: "Check above 7 is not allowed", args: args{flag: 8, filePath: "./test"}, wantErr: true},
		{name: "Check 1 is allowed", args: args{flag: 1, filePath: "./test"}, wantErr: false},
		{name: "Check 7 is allowed", args: args{flag: 7, filePath: "./test"}, wantErr: false},
		{name: "Create a file", args: args{flag: 1, filePath: "./test"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _, err := NewRequestLogger(tt.args.filePath, tt.args.flag)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequestLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.Close()
			os.Remove(tt.args.filePath)
		})
	}
}
