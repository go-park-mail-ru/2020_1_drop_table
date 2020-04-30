package update_functions

import "errors"

var ErrNotInt error

func init() {
	ErrNotInt = errors.New("given value is not int")
}
