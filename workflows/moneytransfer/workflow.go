package moneytransfer

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-samples/workflows/service"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"github.com/indeedeng/iwf-golang-sdk/iwf/ptr"
)

func NewMoneyTransferWorkflow(svc service.MyService) iwf.ObjectWorkflow {

	return &MoneyTransferWorkflow{
		svc: svc,
	}
}

type MoneyTransferWorkflow struct {
	iwf.WorkflowDefaults

	svc service.MyService
}

func (e MoneyTransferWorkflow) GetWorkflowStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(&checkBalanceState{svc: e.svc}),
		iwf.NonStartingStateDef(&createDebitMemoState{svc: e.svc}),
		iwf.NonStartingStateDef(&debitState{svc: e.svc}),
		iwf.NonStartingStateDef(&createCreditMemoState{svc: e.svc}),
		iwf.NonStartingStateDef(&creditState{svc: e.svc}),
		iwf.NonStartingStateDef(&compensateState{svc: e.svc}),
	}
}

type TransferRequest struct {
	FromAccount string
	ToAccount   string
	Amount      int
	Notes       string
}

type checkBalanceState struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (i checkBalanceState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	var request TransferRequest
	input.Get(&request)

	hasSufficientFunds := i.svc.CheckBalance(request.FromAccount, request.Amount)
	if !hasSufficientFunds {
		return iwf.ForceFailWorkflow("insufficient funds"), nil
	}

	return iwf.SingleNextState(&createDebitMemoState{}, request), nil
}

type createDebitMemoState struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (i createDebitMemoState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	var request TransferRequest
	input.Get(&request)

	err := i.svc.CreateDebitMemo(request.FromAccount, request.Amount, request.Notes)
	if err != nil {
		return nil, err
	}

	return iwf.SingleNextState(&debitState{}, request), nil
}

func (i createDebitMemoState) GetStateOptions() *iwf.StateOptions {
	return &iwf.StateOptions{
		ExecuteApiRetryPolicy: &iwfidl.RetryPolicy{
			MaximumAttemptsDurationSeconds: ptr.Any(int32(3600)),
		},
		ExecuteApiFailureProceedState: &compensateState{},
	}
}

type debitState struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (i debitState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	var request TransferRequest
	input.Get(&request)

	err := i.svc.Debit(request.FromAccount, request.Amount)
	if err != nil {
		return nil, err
	}

	return iwf.SingleNextState(&createCreditMemoState{}, request), nil
}

func (i debitState) GetStateOptions() *iwf.StateOptions {
	return &iwf.StateOptions{
		ExecuteApiRetryPolicy: &iwfidl.RetryPolicy{
			MaximumAttemptsDurationSeconds: ptr.Any(int32(3600)),
		},
		ExecuteApiFailureProceedState: &compensateState{},
	}
}

type createCreditMemoState struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (i createCreditMemoState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	var request TransferRequest
	input.Get(&request)

	err := i.svc.CreateCreditMemo(request.ToAccount, request.Amount, request.Notes)
	if err != nil {
		return nil, err
	}

	return iwf.SingleNextState(&creditState{}, request), nil
}

func (i createCreditMemoState) GetStateOptions() *iwf.StateOptions {
	return &iwf.StateOptions{
		ExecuteApiRetryPolicy: &iwfidl.RetryPolicy{
			MaximumAttemptsDurationSeconds: ptr.Any(int32(3600)),
		},
		ExecuteApiFailureProceedState: &compensateState{},
	}
}

type creditState struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (i creditState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	var request TransferRequest
	input.Get(&request)

	err := i.svc.Credit(request.ToAccount, request.Amount)
	if err != nil {
		return nil, err
	}

	return iwf.GracefulCompleteWorkflow(fmt.Sprintf("transfer is done from %v to %v for amount %v", request.FromAccount, request.ToAccount, request.Amount)), nil
}

func (i creditState) GetStateOptions() *iwf.StateOptions {
	return &iwf.StateOptions{
		ExecuteApiRetryPolicy: &iwfidl.RetryPolicy{
			MaximumAttemptsDurationSeconds: ptr.Any(int32(3600)),
		},
		ExecuteApiFailureProceedState: &compensateState{},
	}
}

type compensateState struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (i compensateState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	// NOTE: to improve, we can use iWF data attributes to track whether each step has been attempted to execute
	// and check a flag to see if we should undo it or not

	var request TransferRequest
	input.Get(&request)

	err := i.svc.UndoCredit(request.ToAccount, request.Amount)
	if err != nil {
		return nil, err
	}
	err = i.svc.UndoCreateCreditMemo(request.ToAccount, request.Amount, request.Notes)
	if err != nil {
		return nil, err
	}
	err = i.svc.UndoCreateDebitMemo(request.FromAccount, request.Amount, request.Notes)
	if err != nil {
		return nil, err
	}
	err = i.svc.UndoDebit(request.FromAccount, request.Amount)
	if err != nil {
		return nil, err
	}

	return iwf.GracefulCompleteWorkflow(fmt.Sprintf("transfer is done from %v to %v for amount %v", request.FromAccount, request.ToAccount, request.Amount)), nil
}

func (i compensateState) GetStateOptions() *iwf.StateOptions {
	return &iwf.StateOptions{
		ExecuteApiRetryPolicy: &iwfidl.RetryPolicy{
			MaximumAttemptsDurationSeconds: ptr.Any(int32(86400)),
		},
	}
}
