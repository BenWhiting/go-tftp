package main

import (
	"flag"
	"log"

	"github.com/BenWhiting/go-tftp/internal/constants"
	"github.com/BenWhiting/go-tftp/pkg/tftpd"
)

func main() {
	serverConfig := &tftpd.Config{}
	flag.StringVar(&serverConfig.Address, "h", "127.0.0.1:69", "TFTP server address")
	flag.StringVar(&serverConfig.LogFilePath, "p", "./go-tftp.log", "TFTP server request log file path.")
	flag.IntVar(&serverConfig.FlushInterval, "f", 1, "TFTP server flush period (Seconds).")
	flag.IntVar(&serverConfig.LogFlag, "l", 2, "TFTP server log flag. [1-7].")
	flag.IntVar(&serverConfig.TransferTimeout, "t", 30, "TFTP server transfer timeout (Seconds).")
	flag.IntVar(&serverConfig.RetryTime, "r", 30, "TFTP server transfer retry time (Seconds).")

	flag.Parse()

	serverConfig.AllowedMode = constants.OctetMode

	server, err := tftpd.NewTFTPServer(serverConfig)
	if nil != err {
		log.Fatal(err)
	}

	err = server.Initialize()
	if nil != err {
		log.Fatal(err)
	}

	server.Start()
}
