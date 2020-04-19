package loyalty_systems

type LoyaltySystem interface {
	UpdatingPass(reqLoyaltyInfo, dbLoyaltyInfo string) (newLoyaltyInfo string, err error)
	CreatingCustomer(loyaltyInfo map[string]interface{}) (point, sum int, err error)
	SettingPoints(loyaltyInfo map[string]interface{}, oldPoint, oldSum int) (newPoint, newSum int, err error)
}
