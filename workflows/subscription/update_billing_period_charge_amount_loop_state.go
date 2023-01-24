package subscription

import (
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type updateBillingPeriodChargeAmountLoopState struct {
	iwf.DefaultStateIdAndOptions
}

func (b updateBillingPeriodChargeAmountLoopState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.AllCommandsCompletedRequest(
		iwf.NewSignalCommand("", SignalUpdateBillingPeriodChargeAmount),
	), nil
}

func (b updateBillingPeriodChargeAmountLoopState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var customer Customer
	persistence.GetDataObject(keyCustomer, &customer)

	var newAmount int
	commandResults.GetSignalCommandResultByChannel(SignalUpdateBillingPeriodChargeAmount).SignalValue.Get(&newAmount)

	customer.Subscription.BillingPeriodCharge = newAmount
	persistence.SetDataObject(keyCustomer, customer)

	return iwf.SingleNextState(updateBillingPeriodChargeAmountLoopState{}, nil), nil
}
