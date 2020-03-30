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
	ErrNoPublishedCard = errors.New("no published card for this cafe")
	ErrNoRequestedCard = errors.New("no card with requested params for this cafe")

	//Customer errors
	BadUuid = errors.New("Bad Uuuid")
)
