package tftpd

import (
	"log"
	"net"
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
	activeTransfers map[string]*communication.Transfer
	store           map[string][]byte
	mux             *sync.RWMutex
}

// Config struct for the server
type Config struct {
	AllowedMode   string
	Address       string
	LogFilePath   string
	FlushInterval int
	LogFlag       int
}

// NewTFTPServer creates a new Server instance
func NewTFTPServer(config *Config) (*Server, error) {
	gl, err := logger.NewGeneralLogger(config.LogFlag)
	if nil != err {
		return nil, err
	}

	fl, err := logger.NewRequestLogger(config.LogFilePath, config.LogFlag)
	if nil != err {
		return nil, err
	}

	return &Server{
		Config:       config,
		logger:       gl,
		loggerToFile: fl,
		mux:          &sync.RWMutex{},
	}, nil
}

// Start a Server instance
func (s *Server) Start() error {
	conn, err := net.ListenPacket("udp", s.Config.Address)
	if nil != err {
		return err
	}
	s.Connection = conn

	s.logger.Printf(constants.ServerStartMsg, conn.LocalAddr())

	defer s.Connection.Close()
	s.serve()
	return nil
}

func (s *Server) serve() {
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
	t := time.NewTicker(time.Millisecond * time.Duration(s.Config.FlushInterval))
	for range t.C {
		s.mux.Lock()
		//TODO: FINISH FLUSH
		s.mux.Unlock()
	}
}
