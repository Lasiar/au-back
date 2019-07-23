package base

import (
	"errors"
	"strings"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotForbidden = errors.New("forbidden")
)

func IsBadRequest(err error) bool   { return strings.HasPrefix(err.Error(), ErrBadRequest.Error()) }
func IsUnauthorized(err error) bool { return strings.HasPrefix(err.Error(), ErrUnauthorized.Error()) }
func IsNotForbidden(err error) bool { return strings.HasPrefix(err.Error(), ErrNotForbidden.Error()) }
