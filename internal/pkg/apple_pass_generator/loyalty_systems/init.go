package loyalty_systems

var LoyaltySystems map[string]LoyaltySystem

func init() {
	LoyaltySystems = map[string]LoyaltySystem{
		"coffee_cup": &CoffeeCup{
			InfoVarName:   "cups_count",
			PointsVarName: "coffee_cups",
		},
		"cashback": &CashBack{
			InfoVarName:   "cashback",
			PointsVarName: "points_count",
		},
	}
}
