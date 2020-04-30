package loyalty_systems

import (
	"encoding/json"
	"fmt"
)

type Percents struct {
	purchasesSumVarName string
	discountVarName     string
	newPurchasesVarName string
}

func (p *Percents) UpdatingPass(reqLoyaltyInfo, dbLoyaltyInfo string) (newLoyaltyInfo string, err error) {
	var reqMap map[int]int
	var DBMap map[int]int
	err = UnmarshalEmptyString([]byte(reqLoyaltyInfo), &reqMap)
	if err == nil {
		return reqLoyaltyInfo, nil
	}
	err = UnmarshalEmptyString([]byte(dbLoyaltyInfo), &DBMap)
	if err != nil {
		return "", ErrBadLoyaltyInfo
	}

	return dbLoyaltyInfo, nil
}

func (p *Percents) CreatingCustomer(loyaltyInfo string) (customerPoints, newLoyaltyInfo string,
	err error) {

	return fmt.Sprintf(`{"%s": 0, "%s": 0}`, p.purchasesSumVarName, p.discountVarName),
		loyaltyInfo, nil
}

func (p *Percents) SettingPoints(loyaltyInfo, dbPoints, reqPoints string) (newPoints string, err error) {
	var reqMap map[string]int

	err = UnmarshalEmptyString([]byte(reqPoints), &reqMap)
	if err != nil {
		return "", err
	}

	newPurchase, ok := reqMap[p.newPurchasesVarName]
	if !ok {
		return "", ErrBadReqPoints
	}

	var dbMap map[string]int
	err = UnmarshalEmptyString([]byte(dbPoints), &dbMap)
	if err != nil {
		return "", err
	}

	purchasesSum, ok := dbMap[p.purchasesSumVarName]
	if !ok {
		return "", ErrBadDBPoints
	}

	purchasesSum += newPurchase
	dbMap[p.purchasesSumVarName] = purchasesSum

	var loyaltyMap map[int]int
	err = UnmarshalEmptyString([]byte(loyaltyInfo), &loyaltyMap)
	if err != nil {
		return "", err
	}

	finalDiscount := 0
	for purchasesSumForDiscount := range loyaltyMap {
		if purchasesSum > purchasesSumForDiscount {
			finalDiscount = loyaltyMap[purchasesSumForDiscount]
		}
	}

	dbMap[p.discountVarName] = finalDiscount
	newPointsBytes, err := json.Marshal(dbMap)
	return string(newPointsBytes), err
}
