package loyalty_systems

type CoffeeCup struct {
	VarName string
}

func (c *CoffeeCup) UpdatingPass(reqLoyaltyInfo, dbLoyaltyInfo string) (newLoyaltyInfo string, err error) {
	var reqMap map[string]int
	var DBMap map[string]int

	err = UnmarshalEmptyString([]byte(reqLoyaltyInfo), &reqMap)
	if err != nil {
		return "", err
	}

	err = UnmarshalEmptyString([]byte(dbLoyaltyInfo), &DBMap)
	if err != nil {
		return "", err
	}

	_, reqOk := reqMap[c.VarName]
	_, DBOk := DBMap[c.VarName]

	if reqOk {
		return reqLoyaltyInfo, nil
	} else if DBOk {
		return dbLoyaltyInfo, nil
	}

	return "", ErrBadLoyaltyInfo
}

func (c *CoffeeCup) CreatingCustomer(loyaltyInfo string) (customerPoints, newLoyaltyInfo string,
	err error) {

	return `{"coffee_cups": 0}`, loyaltyInfo, nil
}

func (c *CoffeeCup) SettingPoints(loyaltyInfo, dbPoints, reqPoints string) (newPoints string, err error) {

	return "", nil
}
