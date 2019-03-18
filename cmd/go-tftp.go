package main

import (
	"flag"
	"log"

	"github.com/BenWhiting/go-tftp/internal/constants"
	"github.com/BenWhiting/go-tftp/pkg/tftpd"
)

func main() {
	serverConfig := &tftpd.Config{}

	flag.StringVar(&serverConfig.Address, "h", "127.0.0.1:69", "TFTP server hosting address")
	flag.StringVar(&serverConfig.LogFilePath, "p", "go-tftp.log", "TFTP server request log file path.")
	flag.IntVar(&serverConfig.FlushInterval, "f", 500, "TFTP server flush period in milliseconds.")
	flag.IntVar(&serverConfig.LogFlag, "l", 2, "TFTP server log flag. (1-7)")
	flag.Parse()

	serverConfig.AllowedMode = constants.OctetMode

	server, err := tftpd.NewTFTPServer(serverConfig)
	if nil != err {
		log.Fatal(err)
	}

	err = server.Start()
	if nil != err {
		log.Fatal(err)
	}

}
