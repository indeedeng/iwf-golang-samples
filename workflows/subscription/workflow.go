package subscription

import (
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type SubscriptionWorkflow struct {
	iwf.DefaultWorkflowType
}

const (
	keyBillingPeriodNum = "billingPeriodNum"
	keyCustomer         = "customer"

	SignalCancelSubscription              = "cancelSubscription"
	SignalUpdateBillingPeriodChargeAmount = "updateBillingPeriodChargeAmount"
)

func (b SubscriptionWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(&initState{}),
		iwf.NonStartingStateDef(&cancelState{}),
		iwf.NonStartingStateDef(&updateBillingPeriodChargeAmountLoopState{}),
		iwf.NonStartingStateDef(&trialState{}),
		iwf.NonStartingStateDef(&chargeLoopState{}),
	}
}

func (b SubscriptionWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.DataObjectDef(keyBillingPeriodNum),
		iwf.DataObjectDef(keyCustomer),
	}
}

func (b SubscriptionWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.SignalChannelDef(SignalCancelSubscription),
		iwf.SignalChannelDef(SignalUpdateBillingPeriodChargeAmount),
	}
}
