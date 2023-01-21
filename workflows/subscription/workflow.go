package subscription

import (
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type SubscriptionWorkflow struct{}

const (
	keyBillingPeriodNum = "billingPeriodNum"
	keyCustomer         = "customer"

	signalCancelSubscription              = "cancelSubscription"
	signalUpdateBillingPeriodChargeAmount = "updateBillingPeriodChargeAmount"
)

func (b SubscriptionWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.NewStartingState(&initState{}),
		iwf.NewNonStartingState(&cancelState{}),
		iwf.NewNonStartingState(&updateBillingPeriodChargeAmountLoopState{}),
		iwf.NewNonStartingState(&trialState{}),
		iwf.NewNonStartingState(&chargeLoopState{}),
	}
}

func (b SubscriptionWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.NewDataObjectDef(keyBillingPeriodNum),
		iwf.NewDataObjectDef(keyCustomer),
	}
}

func (b SubscriptionWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.NewSignalChannelDef(signalCancelSubscription),
		iwf.NewSignalChannelDef(signalUpdateBillingPeriodChargeAmount),
	}
}

func (b SubscriptionWorkflow) GetWorkflowType() string {
	return ""
}
