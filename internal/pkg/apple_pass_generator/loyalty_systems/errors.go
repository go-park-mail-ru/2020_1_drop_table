package loyalty_systems

import "errors"

var (
	ErrBadLoyaltyInfo = errors.New("bad loyalty info in database and request")
)
