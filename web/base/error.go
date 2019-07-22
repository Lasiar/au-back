package base

import (
	"errors"
	"strings"
)

var (
	ErrBadRequest = errors.New("bad request")
)

func IsBadRequest(err error) bool { return strings.HasPrefix(err.Error(), ErrBadRequest.Error()) }
