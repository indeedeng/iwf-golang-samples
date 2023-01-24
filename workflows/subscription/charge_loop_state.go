package subscription

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

type chargeLoopState struct {
	iwf.DefaultStateIdAndOptions
}

const subscriptionOverKey = "subscriptionOver"

func (b chargeLoopState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	var customer Customer
	persistence.GetDataObject(keyCustomer, &customer)

	var periodNum int
	persistence.GetDataObject(keyBillingPeriodNum, &periodNum)

	if periodNum >= customer.Subscription.MaxBillingPeriods {
		persistence.SetStateLocal(subscriptionOverKey, true)
		return iwf.EmptyCommandRequest(), nil
	}

	persistence.SetDataObject(keyBillingPeriodNum, periodNum+1)

	return iwf.AllCommandsCompletedRequest(
		iwf.NewTimerCommand("", time.Now().Add(customer.Subscription.BillingPeriod)),
	), nil
}

func (b chargeLoopState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var customer Customer
	persistence.GetDataObject(keyCustomer, &customer)

	var subscriptionOver bool
	persistence.GetStateLocal(subscriptionOverKey, &subscriptionOver)
	if subscriptionOver {
		fmt.Println("this is an RPC call to send a subscription over email to user ", customer.Email)
		// use force completing because the cancel state is still waiting for signal
		return iwf.ForceCompletingWorkflow, nil
	}

	fmt.Printf("this is an RPC call to charge customer %v for $%v \n", customer.Email, customer.Subscription.BillingPeriodCharge)

	return iwf.SingleNextState(chargeLoopState{}, nil), nil
}
