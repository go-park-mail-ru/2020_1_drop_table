package loyalty_systems

import "fmt"

type CoffeeCup struct {
	InfoVarName   string
	PointsVarName string
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

	_, reqOk := reqMap[c.InfoVarName]
	_, DBOk := DBMap[c.InfoVarName]

	if reqOk {
		return reqLoyaltyInfo, nil
	} else if DBOk {
		return dbLoyaltyInfo, nil
	}

	return "", ErrBadLoyaltyInfo
}

func (c *CoffeeCup) CreatingCustomer(loyaltyInfo string) (customerPoints, newLoyaltyInfo string,
	err error) {

	return fmt.Sprintf(`{"%s": 0}`, c.PointsVarName), loyaltyInfo, nil
}

func (c *CoffeeCup) SettingPoints(_, _, reqPoints string) (newPoints string, err error) {
	var reqMap map[string]int

	err = UnmarshalEmptyString([]byte(reqPoints), &reqMap)
	if err != nil {
		return "", err
	}

	pointsReq, reqOk := reqMap[c.PointsVarName]
	if !reqOk {
		return "", ErrBadPoints
	}

	if pointsReq < 0 {
		return "", ErrValidationCoffeeCups
	}

	newPointsJson := fmt.Sprintf(`{"%s": %d}`, c.PointsVarName, pointsReq)
	return newPointsJson, nil
}
