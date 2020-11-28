package paypal

// Plans
type subPlanData struct {
	PlanID          string
	PremiumDuration int
}

// monthlyRegularPlan - planID for 30 day regular price subscription
var monthlyRegularPlan subPlanData = subPlanData{PlanID: "P-96F880834P1216324L7BMCAA", PremiumDuration: 30}
