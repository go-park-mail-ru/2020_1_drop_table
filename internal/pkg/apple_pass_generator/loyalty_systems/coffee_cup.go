package loyalty_systems

type CoffeeCup struct{}

func (c *CoffeeCup) UpdatingPass(_, _ string) (newLoyaltyInfo string, err error) {
	return "{}", nil
}

func (c *CoffeeCup) CreatingCustomer(loyaltyInfo map[string]interface{}) (point, sum int, err error) {
	return 0, 0, nil
}

func (c *CoffeeCup) SettingPoints(loyaltyInfo map[string]interface{}, oldPoint,
	oldSum int) (newPoint, newSum int, err error) {

	return 0, 0, nil
}
