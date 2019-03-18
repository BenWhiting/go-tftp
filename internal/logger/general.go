package logger

import (
	"errors"
	"log"
	"os"

	"github.com/BenWhiting/go-tftp/internal/constants"
)

// NewGeneralLogger for the tftp server
func NewGeneralLogger(flag int) (*log.Logger, error) {
	if flag > 7 || flag < 1 {
		return nil, errors.New(constants.UnknownLogFlagMsg)
	}
	return log.New(os.Stdout, "go-tftp: ", flag), nil
}
