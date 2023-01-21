package subscription

import (
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type SubscriptionWorkflow struct{}

const (
	keySubscriptionCancelled = "subscriptionCancelled"
	keyBillingPeriodNum      = "billingPeriodNum"
	keyCustomer              = "customer"

	signalCancelSubscription = "cancelSubscription"
)

func (b SubscriptionWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.NewStartingState(&trialState{}),
		iwf.NewNonStartingState(&cancelState{}),
	}
}

func (b SubscriptionWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.NewDataObjectDef(keySubscriptionCancelled),
		iwf.NewDataObjectDef(keyBillingPeriodNum),
		iwf.NewDataObjectDef(keyCustomer),
	}
}

func (b SubscriptionWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.NewSignalChannelDef(signalCancelSubscription),
	}
}

func (b SubscriptionWorkflow) GetWorkflowType() string {
	return ""
}
