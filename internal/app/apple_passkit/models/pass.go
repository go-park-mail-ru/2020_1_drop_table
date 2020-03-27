package models

//ToDo write custom validator for Design
type ApplePassDB struct {
	ApplePassID int    `json:"-"`
	Design      string `json:"design" validate:"required"`
	Icon        []byte `json:"icon" validate:"required"`
	Icon2x      []byte `json:"icon2x" validate:"required"`
	Logo        []byte `json:"logo" validate:"required"`
	Logo2x      []byte `json:"logo2x" validate:"required"`
	Strip       []byte `json:"strip" validate:"required"`
	Strip2x     []byte `json:"strip2x" validate:"required"`
}
