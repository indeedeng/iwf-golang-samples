package engagement

import (
	"fmt"
	"github.com/indeedeng/iwf-golang-samples/workflows/service"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

func NewEngagementWorkflow(svc service.MyService) iwf.ObjectWorkflow {

	return &EngagementWorkflow{
		svc: svc,
	}
}

type EngagementWorkflow struct {
	iwf.DefaultWorkflowType

	svc service.MyService
}

func (e EngagementWorkflow) GetWorkflowStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(NewInitState()),
		iwf.NonStartingStateDef(NewProcessTimoutState(e.svc)),
		iwf.NonStartingStateDef(NewReminderState(e.svc)),
		iwf.NonStartingStateDef(NewNotifyExternalSystemState(e.svc)),
	}
}

func (e EngagementWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.SearchAttributeDef(keyEmployerId, iwfidl.KEYWORD),
		iwf.SearchAttributeDef(keyJobSeekerId, iwfidl.KEYWORD),
		iwf.SearchAttributeDef(keyStatus, iwfidl.KEYWORD),
		iwf.SearchAttributeDef(keyLastUpdateTimestamp, iwfidl.INT),

		iwf.DataAttributeDef(keyNotes),
	}
}

func (e EngagementWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.SignalChannelDef(SignalChannelOptOutReminder),
		iwf.InternalChannelDef(InternalChannelCompleteProcess),

		iwf.RPCMethodDef(e.Describe, nil),
		iwf.RPCMethodDef(e.Decline, nil),
		iwf.RPCMethodDef(e.Accept, nil),
	}
}

const (
	keyEmployerId          = "EmployerId"
	keyJobSeekerId         = "JobSeekerId"
	keyStatus              = "EngagementStatus"
	keyLastUpdateTimestamp = "LastUpdateTimeMillis"
	keyNotes               = "notes"

	SignalChannelOptOutReminder    = "OptOutReminder"
	InternalChannelCompleteProcess = "CompleteProcess"
)

func (e EngagementWorkflow) Describe(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (interface{}, error) {

	status := persistence.GetSearchAttributeKeyword(keyStatus)
	employerId := persistence.GetSearchAttributeKeyword(keyEmployerId)
	jobSeekerId := persistence.GetSearchAttributeKeyword(keyJobSeekerId)
	var notes string
	persistence.GetDataAttribute(keyNotes, &notes)

	return EngagementDescription{
		EmployerId:    employerId,
		JobSeekerId:   jobSeekerId,
		Notes:         notes,
		CurrentStatus: Status(status),
	}, nil
}

func (e EngagementWorkflow) Decline(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (interface{}, error) {

	status := Status(persistence.GetSearchAttributeKeyword(keyStatus))
	if status != StatusInitiated {
		return nil, fmt.Errorf("can only decline in INITIATED status, current is %v", status)
	}

	persistence.SetSearchAttributeKeyword(keyStatus, string(StatusDeclined))
	persistence.SetSearchAttributeInt(keyLastUpdateTimestamp, time.Now().Unix())
	communication.TriggerStateMovements(iwf.NewStateMovement(notifyExternalSystemState{}, string(StatusDeclined)))

	var notes string
	input.Get(&notes)

	var currentNotes string
	persistence.GetDataAttribute(keyNotes, &currentNotes)
	persistence.SetDataAttribute(keyNotes, currentNotes+";"+notes)
	return nil, nil
}

func (e EngagementWorkflow) Accept(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (interface{}, error) {

	status := Status(persistence.GetSearchAttributeKeyword(keyStatus))
	if status != StatusInitiated && status != StatusDeclined {
		return nil, fmt.Errorf("can only decline in INITIATED or DECLINED status, current is %v", status)
	}

	persistence.SetSearchAttributeKeyword(keyStatus, string(StatusAccepted))
	persistence.SetSearchAttributeInt(keyLastUpdateTimestamp, time.Now().Unix())
	communication.TriggerStateMovements(iwf.NewStateMovement(notifyExternalSystemState{}, string(StatusAccepted)))

	var notes string
	input.Get(&notes)

	var currentNotes string
	persistence.GetDataAttribute(keyNotes, &currentNotes)
	persistence.SetDataAttribute(keyNotes, currentNotes+";"+notes)
	return nil, nil
}

func NewInitState() iwf.WorkflowState {
	return initState{}
}

type initState struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
}

func (i initState) Execute(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var engInput EngagementInput
	input.Get(&engInput)

	persistence.SetSearchAttributeKeyword(keyEmployerId, engInput.EmployerId)
	persistence.SetSearchAttributeKeyword(keyJobSeekerId, engInput.JobSeekerId)
	persistence.SetSearchAttributeKeyword(keyStatus, string(StatusInitiated))

	persistence.SetDataAttribute(keyNotes, engInput.Notes)
	return iwf.MultiNextStatesWithInput(
		iwf.NewStateMovement(processTimoutState{}, nil),
		iwf.NewStateMovement(reminderState{}, nil),
		iwf.NewStateMovement(notifyExternalSystemState{}, StatusInitiated),
	), nil
}

func NewProcessTimoutState(svc service.MyService) iwf.WorkflowState {
	return processTimoutState{
		svc: svc,
	}
}

type processTimoutState struct {
	iwf.WorkflowStateDefaults
	svc service.MyService
}

func (p processTimoutState) WaitUntil(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.AnyCommandCompletedRequest(
		iwf.NewTimerCommand("", time.Now().Add(time.Hour*24*60)), // ~ 2 months
		iwf.NewInternalChannelCommand("", InternalChannelCompleteProcess),
	), nil
}

func (p processTimoutState) Execute(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	status := persistence.GetSearchAttributeKeyword(keyStatus)
	employerId := persistence.GetSearchAttributeKeyword(keyEmployerId)
	jobSeekerId := persistence.GetSearchAttributeKeyword(keyJobSeekerId)
	updateStatus := "timeout"
	if status == string(StatusAccepted) {
		updateStatus = "done"
	}
	p.svc.UpdateExternalSystem(fmt.Sprintf("notify engagement from employer %v, jobSeeker %v for status %v", employerId, jobSeekerId, status))
	return iwf.GracefulCompleteWorkflow(updateStatus), nil
}

func NewReminderState(svc service.MyService) iwf.WorkflowState {
	return reminderState{
		svc: svc,
	}
}

type reminderState struct {
	iwf.WorkflowStateDefaults
	svc service.MyService
}

func (r reminderState) WaitUntil(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.AnyCommandCompletedRequest(
		iwf.NewTimerCommand("", time.Now().Add(time.Second*5)), // use 5 seconds for demo, should be 24 hours in real world
		iwf.NewSignalCommand("", SignalChannelOptOutReminder),
	), nil
}

func (r reminderState) Execute(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	status := persistence.GetSearchAttributeKeyword(keyStatus)
	if status != string(StatusInitiated) {
		return iwf.DeadEnd, nil
	}
	optoutSignalCommandResult := commandResults.Signals[0]
	if optoutSignalCommandResult.Status == iwfidl.RECEIVED {
		var currentNotes string
		persistence.GetDataAttribute(keyNotes, &currentNotes)
		persistence.SetDataAttribute(keyNotes, currentNotes+";"+"User optout reminder")

		return iwf.DeadEnd, nil
	}

	jobSeekerId := persistence.GetSearchAttributeKeyword(keyJobSeekerId)
	r.svc.SendEmail(jobSeekerId, "Reminder:xxx please respond", "Hello xxx, ...")
	return iwf.SingleNextState(reminderState{}, nil), nil
}

func NewNotifyExternalSystemState(svc service.MyService) iwf.WorkflowState {
	return notifyExternalSystemState{
		svc: svc,
	}
}

type notifyExternalSystemState struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (n notifyExternalSystemState) Execute(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var status Status
	input.Get(&status)

	jobSeekerId := persistence.GetSearchAttributeKeyword(keyJobSeekerId)
	employerId := persistence.GetSearchAttributeKeyword(keyEmployerId)
	n.svc.UpdateExternalSystem(fmt.Sprintf("notify engagement from employerId %v to jobSeekerId %v for status %v ", employerId, jobSeekerId, status))
	return iwf.DeadEnd, nil
}

// GetStateOptions customize the state options
// By default, all state execution will retry infinitely (until workflow timeout).
// This may not work for some dependency as we may want to retry for only a certain times
func (n notifyExternalSystemState) GetStateOptions() *iwf.StateOptions {
	return &iwf.StateOptions{
		ExecuteApiRetryPolicy: &iwfidl.RetryPolicy{
			BackoffCoefficient:             iwfidl.PtrFloat32(2),
			MaximumAttempts:                iwfidl.PtrInt32(100),
			MaximumAttemptsDurationSeconds: iwfidl.PtrInt32(3600),
			MaximumIntervalSeconds:         iwfidl.PtrInt32(60),
			InitialIntervalSeconds:         iwfidl.PtrInt32(3),
		},
	}
}
