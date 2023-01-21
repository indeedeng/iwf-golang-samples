package subscription

import (
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type initState struct{}

const initStateId = "initState"

func (b initState) GetStateId() string {
	return initStateId
}

func (b initState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	var customer Customer
	err := input.Get(&customer)
	if err != nil {
		return nil, err
	}
	err = persistence.SetDataObject(keyCustomer, customer)
	if err != nil {
		return nil, err
	}
	return iwf.EmptyCommandRequest(), nil
}

func (b initState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	return iwf.MultiNextStatesByStateIds(trialStateId, cancelStateId, updateBillingPeriodChargeAmountLoopStateId), nil
}

func (b initState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
