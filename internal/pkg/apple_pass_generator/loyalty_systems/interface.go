package loyalty_systems

type LoyaltySystem interface {
	UpdatingPass(reqLoyaltyInfo, dbLoyaltyInfo string) (newLoyaltyInfo string, err error)
	CreatingCustomer(loyaltyInfo string) (customerPoints string, newLoyaltyInfo string, err error)
	SettingPoints(loyaltyInfo, dbPoints, reqPoints string) (newPoints string, err error)
}
