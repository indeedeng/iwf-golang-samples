package subscription

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

type trialState struct{}

const TrialStateId = "trialState"

func (b trialState) GetStateId() string {
	return TrialStateId
}

func (b trialState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	var customer Customer
	err := input.Get(&customer)
	if err != nil {
		return nil, err
	}
	err = persistence.SetDataObject(keyCustomer, customer)
	if err != nil {
		return nil, err
	}
	err = persistence.SetDataObject(keySubscriptionCancelled, false)
	if err != nil {
		return nil, err
	}

	// send welcome email
	fmt.Println("this is an RPC call to send an welcome email to ", customer.FirstName, customer.LastName, customer.Email)

	return iwf.AllCommandsCompletedRequest(
		iwf.NewTimerCommand("", time.Now().Add(customer.Subscription.TrialPeriod)),
	), nil
}

func (b trialState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var customer Customer
	err := input.Get(&customer)
	if err != nil {
		return nil, err
	}

	var subCanceled bool
	err = persistence.GetDataObject(keySubscriptionCancelled, &subCanceled)
	if err != nil {
		return nil, err
	}
	if subCanceled {
		// send cancel email
		fmt.Println("this is an RPC call to send an cancel email to ", customer.FirstName, customer.LastName, customer.Email)
	}

	err = persistence.SetDataObject(keyBillingPeriodNum, 0)
	if err != nil {
		return nil, err
	}
	return iwf.SingleNextState(ChargeLoopStateId, nil), nil
}

func (b trialState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
