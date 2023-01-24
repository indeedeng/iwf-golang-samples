package subscription

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

type trialState struct {
	iwf.DefaultStateIdAndOptions
}

func (b trialState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	var customer Customer
	persistence.GetDataObject(keyCustomer, &customer)

	// send welcome email
	fmt.Println("this is an RPC call to send an welcome email to ", customer.FirstName, customer.LastName, customer.Email)

	return iwf.AllCommandsCompletedRequest(
		iwf.NewTimerCommand("", time.Now().Add(customer.Subscription.TrialPeriod)),
	), nil
}

func (b trialState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	persistence.SetDataObject(keyBillingPeriodNum, 0)
	return iwf.SingleNextState(chargeLoopState{}, nil), nil
}
