package paypal

// Plans
type subPlanData struct {
	PlanID          string
	PremiumDuration int
}

// monthlyRegularPlan - planID for 30 day regular price subscription
var monthlyRegularPlan subPlanData = subPlanData{PlanID: "P-7X366338T15169429L67OAFQ", PremiumDuration: 30}
