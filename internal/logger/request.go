package logger

import (
	"errors"
	"log"
	"os"

	"github.com/BenWhiting/go-tftp/internal/constants"
)

// NewRequestLogger for tftp server
func NewRequestLogger(filePath string, flag int) (*os.File, *log.Logger, error) {
	if flag > 7 || flag < 1 {
		return nil, nil, errors.New(constants.UnknownLogFlagMsg)
	}
	lf, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		return nil, nil, err
	}
	return lf, log.New(lf, constants.FileRequestMsg, flag), nil
}
