package db

import (
	"github.com/pomkac/mnemonic"
	"errors"
)

var (
	errorInvalidId    = errors.New("Invalid ID")
	errorInvalidValue = errors.New("Invalid value")

	DB *mnemonic.Database
)

func init() {
	config := &mnemonic.Config{ThreadPoolSize: 100,
		ConnPoolSize: 100}
	DB = mnemonic.NewDB(config)
}
