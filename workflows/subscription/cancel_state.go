package subscription

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type cancelState struct{}

const CancelStateId = "cancelState"

func (b cancelState) GetStateId() string {
	return CancelStateId
}

func (b cancelState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.AllCommandsCompletedRequest(
		iwf.NewSignalCommand("", signalCancelSubscription),
	), nil
}

func (b cancelState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var customer Customer
	err := input.Get(&customer)
	if err != nil {
		return nil, err
	}

	err = persistence.SetDataObject(keySubscriptionCancelled, true)
	if err != nil {
		return nil, err
	}

	fmt.Println("this is an RPC call to send CancellationEmailDuringActiveSubscription", customer.Email)
	return iwf.DeadEnd, nil
}

func (b cancelState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
