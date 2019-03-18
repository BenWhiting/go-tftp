package logger

import (
	"log"
	"os"

	"github.com/BenWhiting/go-tftp/internal/constants"
)

// NewRequestLogger for tftp server
func NewRequestLogger(filePath string, flag int) (*log.Logger, error) {
	lf, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		return nil, err
	}
	return log.New(lf, constants.FileRequestMsg, flag), nil
}
