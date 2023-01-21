package subscription

import (
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type updateBillingPeriodChargeAmountLoopState struct{}

const updateBillingPeriodChargeAmountLoopStateId = "updateBillingPeriodChargeAmountLoopState"

func (b updateBillingPeriodChargeAmountLoopState) GetStateId() string {
	return updateBillingPeriodChargeAmountLoopStateId
}

func (b updateBillingPeriodChargeAmountLoopState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.AllCommandsCompletedRequest(
		iwf.NewSignalCommand("", SignalUpdateBillingPeriodChargeAmount),
	), nil
}

func (b updateBillingPeriodChargeAmountLoopState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var customer Customer
	err := persistence.GetDataObject(keyCustomer, &customer)
	if err != nil {
		return nil, err
	}

	var newAmount int
	err = commandResults.GetSignalCommandResultByChannel(SignalUpdateBillingPeriodChargeAmount).SignalValue.Get(&newAmount)
	if err != nil {
		return nil, err
	}

	customer.Subscription.BillingPeriodCharge = newAmount
	err = persistence.SetDataObject(keyCustomer, customer)
	if err != nil {
		return nil, err
	}

	return iwf.SingleNextState(updateBillingPeriodChargeAmountLoopStateId, nil), nil
}

func (b updateBillingPeriodChargeAmountLoopState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
