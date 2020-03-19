package models

import "errors"

var (
	ErrNotFound = errors.New("resource you request not found")
	ErrExisted  = errors.New("given item already existed")

	ErrForbidden = errors.New("no permission")

	ErrBadRequest = errors.New("bad request")

	ErrEmptyJSON = errors.New("empty jsonData field")
	ErrBadJSON   = errors.New("json parsing error")
)
