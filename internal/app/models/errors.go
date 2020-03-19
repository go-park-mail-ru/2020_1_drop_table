package models

import "errors"

var (
	ErrNotFound = errors.New("resource you request not found")
	ErrExisted  = errors.New("given item already existed")

	ErrForbidden = errors.New("no permission")

	ErrInvalidAction = errors.New("you're trying to edit not your cafe")

	ErrBadRequest = errors.New("bad request")

	ErrEmptyJSON = errors.New("empty jsonData field")
	ErrBadJSON   = errors.New("json parsing error")
)
