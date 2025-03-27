package dto

import "errors"

var ErrDuplicated = errors.New("duplicated")
var ErrItemNotFound = errors.New("not_found")
var ErrGameAlreadyStarted = errors.New("game_started")
