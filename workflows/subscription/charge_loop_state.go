package subscription

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

type chargeLoopState struct{}

const ChargeLoopStateId = "chargeLoopState"

const subscriptionOverKey = "subscriptionOver"

func (b chargeLoopState) GetStateId() string {
	return ChargeLoopStateId
}

func (b chargeLoopState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	var customer Customer
	err := input.Get(&customer)
	if err != nil {
		return nil, err
	}

	var periodNum int
	err = persistence.GetDataObject(keyBillingPeriodNum, &periodNum)
	if err != nil {
		return nil, err
	}

	if periodNum >= customer.Subscription.MaxBillingPeriods {
		err := persistence.SetStateLocal(subscriptionOverKey, true)
		if err != nil {
			return nil, err
		}
		return iwf.EmptyCommandRequest(), nil
	}

	err = persistence.SetDataObject(keyBillingPeriodNum, periodNum+1)
	if err != nil {
		return nil, err
	}

	return iwf.AllCommandsCompletedRequest(
		iwf.NewTimerCommand("", time.Now().Add(customer.Subscription.BillingPeriod)),
	), nil
}

func (b chargeLoopState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var customer Customer
	err := input.Get(&customer)
	if err != nil {
		return nil, err
	}

	var subscriptionOver bool
	err = persistence.GetStateLocal(subscriptionOverKey, &subscriptionOver)
	if err != nil {
		return nil, err
	}
	if subscriptionOver {
		fmt.Println("this is an RPC call to send a subscription over email to user ", customer.Email)
		// use force completing because the cancel state is still waiting for signal
		return iwf.ForceCompletingWorkflow, nil
	}

	fmt.Printf("this is an RPC call to charge customer %v for $%v \n", customer.Email, customer.Subscription.BillingPeriodCharge)

	return iwf.SingleNextState(ChargeLoopStateId, nil), nil
}

func (b chargeLoopState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
