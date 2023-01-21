package subscription

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type cancelState struct{}

const cancelStateId = "cancelState"

func (b cancelState) GetStateId() string {
	return cancelStateId
}

func (b cancelState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.AllCommandsCompletedRequest(
		iwf.NewSignalCommand("", signalCancelSubscription),
	), nil
}

func (b cancelState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var customer Customer
	err := persistence.GetDataObject(keyCustomer, &customer)
	if err != nil {
		return nil, err
	}

	fmt.Println("this is an RPC call to send a cancellation email", customer.Email)
	return iwf.ForceCompletingWorkflow, nil
}

func (b cancelState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
