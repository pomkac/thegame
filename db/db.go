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
	config := &mnemonic.Config{ThreadPoolSize: 5000,
		ConnPoolSize: 5000}
	DB = mnemonic.NewDB(config)
}
