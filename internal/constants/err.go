package constants

import "errors"

var (
	ErrNotToPerson   = errors.New("no idOfWallet to person")
	ErrNotFromPerson = errors.New("no idOfWallet from person")
)
