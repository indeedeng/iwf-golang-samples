package subscription

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

type SubscriptionWorkflow struct {
	iwf.DefaultWorkflowType

	svc MyService
}

func NewSubscriptionWorkflow(svc MyService) *SubscriptionWorkflow {
	return &SubscriptionWorkflow{
		svc: svc,
	}
}

const (
	keyBillingPeriodNum = "billingPeriodNum"
	keyCustomer         = "customer"

	SignalCancelSubscription              = "cancelSubscription"
	SignalUpdateBillingPeriodChargeAmount = "updateBillingPeriodChargeAmount"
)

func (b SubscriptionWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(NewInitState()),
		iwf.NonStartingStateDef(NewTriageState(b.svc)),
		iwf.NonStartingStateDef(NewChargeCurrentBillState(b.svc)),
		iwf.NonStartingStateDef(NewCancelState(b.svc)),
		iwf.NonStartingStateDef(NewUpdateChargeAmountState()),
	}
}

func (b SubscriptionWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.DataObjectDef(keyBillingPeriodNum),
		iwf.DataObjectDef(keyCustomer),
	}
}

func (b SubscriptionWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.SignalChannelDef(SignalCancelSubscription),
		iwf.SignalChannelDef(SignalUpdateBillingPeriodChargeAmount),
	}
}

type Subscription struct {
	TrialPeriod         time.Duration
	BillingPeriod       time.Duration
	MaxBillingPeriods   int
	BillingPeriodCharge int
}

type Customer struct {
	FirstName    string
	LastName     string
	Id           string
	Email        string
	Subscription Subscription
}

func NewInitState() iwf.WorkflowState {
	return initState{}
}

type initState struct {
	iwf.DefaultStateIdAndOptions
}

func (b initState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	var customer Customer
	input.Get(&customer)
	persistence.SetDataObject(keyCustomer, customer)
	return iwf.EmptyCommandRequest(), nil
}

func (b initState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	return iwf.MultiNextStates(trialState{}, cancelState{}, updateChargeAmountState{}), nil
}

func NewTriageState(svc MyService) iwf.WorkflowState {
	return trialState{
		svc: svc,
	}
}

type trialState struct {
	iwf.DefaultStateIdAndOptions
	svc MyService
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
	return iwf.SingleNextState(chargeCurrentBillState{}, nil), nil
}

func NewChargeCurrentBillState(svc MyService) iwf.WorkflowState {
	return chargeCurrentBillState{
		svc: svc,
	}
}

type chargeCurrentBillState struct {
	iwf.DefaultStateIdAndOptions
	svc MyService
}

const subscriptionOverKey = "subscriptionOver"

func (b chargeCurrentBillState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
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

func (b chargeCurrentBillState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
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

	return iwf.SingleNextState(chargeCurrentBillState{}, nil), nil
}

func NewCancelState(svc MyService) iwf.WorkflowState {
	return cancelState{
		svc: svc,
	}
}

type cancelState struct {
	iwf.DefaultStateIdAndOptions
	svc MyService
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

func NewUpdateChargeAmountState() iwf.WorkflowState {
	return updateChargeAmountState{}
}

type updateChargeAmountState struct {
	iwf.DefaultStateIdAndOptions
}

func (b updateChargeAmountState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.AllCommandsCompletedRequest(
		iwf.NewSignalCommand("", SignalUpdateBillingPeriodChargeAmount),
	), nil
}

func (b updateChargeAmountState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var customer Customer
	persistence.GetDataObject(keyCustomer, &customer)

	var newAmount int
	commandResults.GetSignalCommandResultByChannel(SignalUpdateBillingPeriodChargeAmount).SignalValue.Get(&newAmount)

	customer.Subscription.BillingPeriodCharge = newAmount
	persistence.SetDataObject(keyCustomer, customer)

	return iwf.SingleNextState(updateChargeAmountState{}, nil), nil
}
