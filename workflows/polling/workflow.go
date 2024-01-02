package polling

import (
	"github.com/indeedeng/iwf-golang-samples/workflows/service"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

func NewPollingWorkflow(svc service.MyService) iwf.ObjectWorkflow {

	return &PollingWorkflow{
		svc: svc,
	}
}

const (
	dataAttrCurrPolls = "currPolls" // tracks how many polls have been done

	SignalChannelTaskACompleted = "taskACompleted"
	SignalChannelTaskBCompleted = "taskBCompleted"

	InternalChannelTaskCCompleted = "taskCCompleted"
)

type PollingWorkflow struct {
	iwf.WorkflowDefaults

	svc service.MyService
}

func (e PollingWorkflow) GetWorkflowStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(&initState{}),
		iwf.NonStartingStateDef(&pollState{svc: e.svc}),
		iwf.NonStartingStateDef(&checkAndCompleteState{svc: e.svc}),
	}
}

func (e PollingWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.DataAttributeDef(dataAttrCurrPolls),
	}
}

func (e PollingWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.SignalChannelDef(SignalChannelTaskACompleted),
		iwf.SignalChannelDef(SignalChannelTaskBCompleted),
		iwf.InternalChannelDef(InternalChannelTaskCCompleted),
	}
}

type initState struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
}

func (i initState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	var maxPollsRequired int
	input.Get(&maxPollsRequired)

	return iwf.MultiNextStatesWithInput(
		iwf.NewStateMovement(pollState{}, maxPollsRequired),
		iwf.NewStateMovement(checkAndCompleteState{}, nil),
	), nil
}

type checkAndCompleteState struct {
	iwf.WorkflowStateDefaults
	svc service.MyService
}

func (i checkAndCompleteState) WaitUntil(
	ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication,
) (*iwf.CommandRequest, error) {
	return iwf.AllCommandsCompletedRequest(
		iwf.NewSignalCommand("", SignalChannelTaskACompleted),
		iwf.NewSignalCommand("", SignalChannelTaskBCompleted),
		iwf.NewInternalChannelCommand("", InternalChannelTaskCCompleted),
	), nil
}

func (i checkAndCompleteState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	return iwf.GracefulCompletingWorkflow, nil
}

type pollState struct {
	iwf.WorkflowStateDefaults
	svc service.MyService
}

func (i pollState) WaitUntil(
	ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication,
) (*iwf.CommandRequest, error) {

	return iwf.AnyCommandCompletedRequest(
		iwf.NewTimerCommand("", time.Now().Add(time.Second*2)),
	), nil
}

func (i pollState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	var maxPollsRequired int
	input.Get(&maxPollsRequired)

	i.svc.CallAPI1("calling API1 for polling service C")

	var currPolls int
	persistence.GetDataAttribute(dataAttrCurrPolls, &currPolls)
	if currPolls >= maxPollsRequired {
		communication.PublishInternalChannel(InternalChannelTaskCCompleted, nil)
		return iwf.DeadEnd, nil
	}

	persistence.SetDataAttribute(dataAttrCurrPolls, currPolls+1)
	// loop back to check
	return iwf.SingleNextState(pollState{}, maxPollsRequired), nil
}