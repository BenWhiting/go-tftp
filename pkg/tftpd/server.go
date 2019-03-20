package tftpd

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/BenWhiting/go-tftp/internal/constants"
	"github.com/BenWhiting/go-tftp/internal/logger"
	"github.com/BenWhiting/go-tftp/pkg/tftpd/communication"
	"github.com/BenWhiting/go-tftp/pkg/wire"
)

// Server struct defining the TFTP deployment
type Server struct {
	Connection      net.PacketConn
	Config          *Config
	logger          *log.Logger
	loggerToFile    *log.Logger
	logFile         *os.File
	activeTransfers map[string]*communication.Transfer
	store           map[string][]byte
	mux             *sync.RWMutex
}

// Config struct for the server
type Config struct {
	AllowedMode     string
	Address         string
	LogFilePath     string
	FlushInterval   int
	LogFlag         int
	TransferTimeout int
	RetryTime       int
}

// NewTFTPServer creates a new Server instance
func NewTFTPServer(config *Config) (*Server, error) {
	// force no values below 0
	if config.FlushInterval <= 0 {
		return nil, fmt.Errorf(constants.ConfigErrMsg, "FlushInterval")
	}
	if config.RetryTime <= 0 {
		return nil, fmt.Errorf(constants.ConfigErrMsg, "RetryTime")
	}
	if config.TransferTimeout <= 0 {
		return nil, fmt.Errorf(constants.ConfigErrMsg, "TransferTimeout")
	}

	gl, err := logger.NewGeneralLogger(config.LogFlag)
	if nil != err {
		return nil, err
	}

	f, fl, err := logger.NewRequestLogger(config.LogFilePath, config.LogFlag)
	if nil != err {
		return nil, err
	}

	return &Server{
		Config:          config,
		logger:          gl,
		loggerToFile:    fl,
		logFile:         f,
		activeTransfers: map[string]*communication.Transfer{},
		store:           map[string][]byte{},
		mux:             &sync.RWMutex{},
	}, nil
}

// Initialize a Server instance
func (s *Server) Initialize() error {
	conn, err := net.ListenPacket("udp", s.Config.Address)
	if nil != err {
		return err
	}
	s.Connection = conn

	s.logger.Printf(constants.ServerStartMsg, conn.LocalAddr())
	return nil
}

// Start the server
func (s *Server) Start() {
	defer s.Connection.Close()
	defer s.logFile.Close()
	// start "pipe" cleaner
	go s.flush()
	// blocker call
	s.readAndRoute()
}

func (s *Server) readAndRoute() {
	for {
		buffer := make([]byte, wire.MaxPacketSize)
		n, addr, err := s.Connection.ReadFrom(buffer)
		if nil != err {
			s.logger.Printf(constants.SimpleErrMsg, err)
			continue
		}
		go s.route(s.Connection, addr, buffer, n)
	}
}

func (s *Server) flush() {
	t := time.NewTicker(time.Second * time.Duration(s.Config.FlushInterval))
	for range t.C {
		s.mux.Lock()
		for _, t := range s.activeTransfers {
			td := time.Duration(time.Second * time.Duration(s.Config.TransferTimeout))
			if td < time.Now().Sub(t.LastOp) {
				// Timed out
				s.logger.Printf(constants.TransferTimeoutMsg, t.Filename)
				delete(s.activeTransfers, t.Filename)
			} else if t.Retry {
				// retry
				rd := time.Duration(time.Second * time.Duration(s.Config.RetryTime))
				if rd < time.Now().Sub(t.LastOp) {
					s.logger.Printf(constants.RetryLastPktMsg, t.Filename)
					t.Transmit()
				}
			}
		}
		s.mux.Unlock()
	}
}
