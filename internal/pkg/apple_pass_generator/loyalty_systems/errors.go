package loyalty_systems

import "errors"

var (
	ErrBadLoyaltyInfo = errors.New("bad loyalty info in database and request")

	ErrBadPoints    = errors.New("bad points in request or DB")
	ErrBadReqPoints = errors.New("bad points in request")
	ErrBadDBPoints  = errors.New("bad points in DB")

	ErrValidationCoffeeCups = errors.New("points count have to be more then 0")
)
