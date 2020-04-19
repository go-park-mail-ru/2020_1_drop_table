package models

import "errors"

var (
	ErrNotFound = errors.New("resource you request not found")
	ErrExisted  = errors.New("given item already existed")

	ErrForbidden = errors.New("no permission")

	ErrInvalidAction = errors.New("you're trying to edit not your cafe")

	ErrBadRequest = errors.New("bad request")

	ErrBadURLParams = errors.New("bad params in URL")

	//Request files errors
	ErrUnexpectedFile         = errors.New("too many files(>1) in one field")
	ErrUnexpectedFilenameText = "unexpected filename: %s"

	//JSON errors
	ErrEmptyJSON = errors.New("empty jsonData field")
	ErrBadJSON   = errors.New("json parsing error")

	//Apple pass errors
	ErrNoLoyaltyProgram = errors.New("no loyalty program with given name")

	//Customer errors
	ErrBadUuid     = errors.New("bad uuid")
	ErrPointsError = errors.New("points can't be <0")

	//Auth errors
	ErrUnauthorized = errors.New("incorrect password or email")
)
