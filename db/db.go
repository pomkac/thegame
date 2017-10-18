package db

import (
	"github.com/pomkac/mnemonic"
	"errors"
)

var (
	errorInvalidId    = errors.New("Invalid ID")
	errorInvalidValue = errors.New("Invalid value")

	DB = mnemonic.NewDB()
)
