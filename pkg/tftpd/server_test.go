package tftpd

import (
	"fmt"
	"os"
	"testing"
)

func TestNewTFTPServer(t *testing.T) {
	validConfig := &Config{
		AllowedMode:     "octet",
		Address:         "127.0.0.1:69",
		LogFilePath:     "./test.txt",
		FlushInterval:   1,
		LogFlag:         1,
		TransferTimeout: 1,
		RetryTime:       1,
	}
	invalidConfig := &Config{
		AllowedMode:     "",
		Address:         "",
		LogFilePath:     "./test.txt",
		FlushInterval:   1,
		LogFlag:         0,
		TransferTimeout: 1,
		RetryTime:       1,
	}
	invalidLogPathConfig := &Config{
		AllowedMode:     "octet",
		Address:         "127.0.0.1:69",
		LogFilePath:     "./?/.&6#$&02/../test.txt",
		FlushInterval:   1,
		LogFlag:         1,
		TransferTimeout: 1,
		RetryTime:       1,
	}
	lowTransferTimeoutConfig := &Config{
		AllowedMode:     "A FAKE MODE",
		Address:         "FOO:BAR",
		LogFilePath:     "./test.txt",
		FlushInterval:   1,
		LogFlag:         1,
		TransferTimeout: 0,
		RetryTime:       1,
	}
	lowRetryTimeConfig := &Config{
		AllowedMode:     "A FAKE MODE",
		Address:         "FOO:BAR",
		LogFilePath:     "./test.txt",
		FlushInterval:   1,
		LogFlag:         1,
		TransferTimeout: 1,
		RetryTime:       0,
	}
	lowFlushIntervalConfig := &Config{
		AllowedMode:     "A FAKE MODE",
		Address:         "FOO:BAR",
		LogFilePath:     "./test.txt",
		FlushInterval:   0,
		LogFlag:         1,
		TransferTimeout: 1,
		RetryTime:       1,
	}

	type args struct {
		config *Config
	}
	tests := []struct {
		name        string
		args        args
		want        *Server
		checkConfig bool
		clearFile   bool
		wantErr     bool
	}{
		{name: "Valid Config Presented", args: args{config: validConfig}, wantErr: false, checkConfig: true, clearFile: true},
		{name: "Invalid log flag", args: args{config: invalidConfig}, wantErr: true},
		{name: "Invalid config file log path", args: args{config: invalidLogPathConfig}, wantErr: true},
		{name: "Invalid low transfer timeout", args: args{config: lowTransferTimeoutConfig}, wantErr: true, checkConfig: false},
		{name: "Invalid low retry timeout", args: args{config: lowRetryTimeConfig}, wantErr: true, checkConfig: false},
		{name: "Invalid low flush interval", args: args{config: lowFlushIntervalConfig}, wantErr: true, checkConfig: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTFTPServer(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTFTPServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkConfig {
				if got.Config.Address != tt.args.config.Address {
					t.Errorf("NewTFTPServer() config Address lost")
					return
				}

				if got.Config.AllowedMode != tt.args.config.AllowedMode {
					t.Errorf("NewTFTPServer() config AllowedMode lost")
					return
				}

				if got.Config.FlushInterval != tt.args.config.FlushInterval {
					t.Errorf("NewTFTPServer() config FlushInterval lost")
					return
				}

				if got.Config.LogFilePath != tt.args.config.LogFilePath {
					t.Errorf("NewTFTPServer() config LogFilePath lost")
					return
				}

				if got.Config.LogFlag != tt.args.config.LogFlag {
					t.Errorf("NewTFTPServer() config LogFlag lost")
					return
				}

				if got.Config.RetryTime != tt.args.config.RetryTime {
					t.Errorf("NewTFTPServer() config RetryTime lost")
					return
				}

				if got.Config.TransferTimeout != tt.args.config.TransferTimeout {
					t.Errorf("NewTFTPServer() config TransferTimeout lost")
					return
				}
			}
			if tt.clearFile {
				got.logFile.Close()
				os.Remove(got.Config.LogFilePath)
			}
		})
	}
}

func TestServer_Initialize(t *testing.T) {
	validConfig := &Config{
		AllowedMode:     "octet",
		Address:         "127.0.0.1:69",
		LogFilePath:     "./test.txt",
		FlushInterval:   30,
		LogFlag:         5,
		TransferTimeout: 30,
		RetryTime:       30,
	}
	invalidAddressConfig := &Config{
		AllowedMode:     "octet",
		Address:         "asdasdasdasd",
		LogFilePath:     "./test.txt",
		FlushInterval:   1,
		LogFlag:         1,
		TransferTimeout: 1,
		RetryTime:       1,
	}
	tests := []struct {
		name      string
		c         *Config
		clearFile bool
		wantErr   bool
	}{
		{name: "Start Valid Server", c: validConfig, wantErr: false, clearFile: true},
		{name: "Bad ADdress for Server", c: invalidAddressConfig, wantErr: true, clearFile: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewTFTPServer(tt.c)
			if nil != err {
				t.Errorf("Set up failure %v", err)
				return
			}
			if err := s.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("Server.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.clearFile {
				s.logFile.Close()
				err := os.Remove(s.Config.LogFilePath)
				if nil != err {
					fmt.Println(err)
				}
			}
		})
	}
}
