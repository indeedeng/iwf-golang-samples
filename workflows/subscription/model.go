package subscription

import "time"

type Subscription struct {
	TrialPeriod         time.Duration
	BillingPeriod       time.Duration
	MaxBillingPeriods   int
	BillingPeriodCharge int
}

type Customer struct {
	FirstName    string
	LastName     string
	Id           string
	Email        string
	Subscription Subscription
}