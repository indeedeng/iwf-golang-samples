package subscription

import (
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type initState struct {
	iwf.DefaultStateIdAndOptions
}

func (b initState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	var customer Customer
	input.Get(&customer)
	persistence.SetDataObject(keyCustomer, customer)
	return iwf.EmptyCommandRequest(), nil
}

func (b initState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	return iwf.MultiNextStates(trialState{}, cancelState{}, updateBillingPeriodChargeAmountLoopState{}), nil
}
