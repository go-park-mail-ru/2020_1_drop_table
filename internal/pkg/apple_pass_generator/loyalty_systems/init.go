package loyalty_systems

var LoyaltySystems map[string]LoyaltySystem

func init() {
	LoyaltySystems = map[string]LoyaltySystem{
		"coffee_cup": &CoffeeCup{VarName: "cups_count"},
	}
}
