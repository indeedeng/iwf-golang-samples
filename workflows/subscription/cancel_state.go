package subscription

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type cancelState struct {
	iwf.DefaultStateIdAndOptions
}

func (b cancelState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.AllCommandsCompletedRequest(
		iwf.NewSignalCommand("", SignalCancelSubscription),
	), nil
}

func (b cancelState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var customer Customer
	persistence.GetDataObject(keyCustomer, &customer)

	fmt.Println("this is an RPC call to send a cancellation email", customer.Email)
	return iwf.ForceCompletingWorkflow, nil
}
