package models

type ApplePassDB struct {
	ApplePassID int    `json:"-"`
	CafeID      int    `json:"cafe_id"`
	Type        string `json:"type"`
	LoyaltyInfo string `json:"loyalty_info" faker:"-"`
	Published   bool   `json:"published"`
	Design      string `json:"design" validate:"required"`
	Icon        []byte `json:"icon" validate:"required"`
	Icon2x      []byte `json:"icon2x" validate:"required"`
	Logo        []byte `json:"logo" validate:"required"`
	Logo2x      []byte `json:"logo2x" validate:"required"`
	Strip       []byte `json:"strip"`
	Strip2x     []byte `json:"strip2x"`
}

type ApplePassMeta struct {
	CafeID int                    `structs:"-"`
	Meta   map[string]interface{} `structs:"Meta"`
}

type UpdateResponse struct {
	URL string
	QR  string
}
