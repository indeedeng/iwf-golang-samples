package microservices

import (
	"github.com/indeedeng/iwf-golang-samples/workflows/service"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

func NewMicroserviceOrchestrationWorkflow(svc service.MyService) iwf.ObjectWorkflow {

	return &OrchestrationWorkflow{
		svc: svc,
	}
}

type OrchestrationWorkflow struct {
	iwf.DefaultWorkflowType

	svc service.MyService
}

func (e OrchestrationWorkflow) GetWorkflowStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(NewState1(e.svc)),
		iwf.NonStartingStateDef(NewState2(e.svc)),
		iwf.NonStartingStateDef(NewState3(e.svc)),
		iwf.NonStartingStateDef(NewState4(e.svc)),
	}
}

func (e OrchestrationWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.DataAttributeDef(keyData),
	}
}

func (e OrchestrationWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.SignalChannelDef(SignalChannelReady),

		iwf.RPCMethodDef(e.Swap, nil),
	}
}

const (
	keyData = "data"

	SignalChannelReady = "Ready"
)

func (e OrchestrationWorkflow) Swap(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (interface{}, error) {

	var oldData string
	persistence.GetDataAttribute(keyData, &oldData)
	var newData string
	input.Get(&newData)
	persistence.SetDataAttribute(keyData, newData)

	return oldData, nil
}

func NewState1(svc service.MyService) iwf.WorkflowState {
	return state1{svc: svc}
}

type state1 struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (i state1) Execute(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var inString string
	input.Get(&inString)

	i.svc.CallAPI1(inString)

	persistence.SetDataAttribute(keyData, inString)
	return iwf.MultiNextStatesWithInput(
		iwf.NewStateMovement(state2{}, nil),
		iwf.NewStateMovement(state3{}, nil),
	), nil
}

func NewState2(svc service.MyService) iwf.WorkflowState {
	return state2{svc: svc}
}

type state2 struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (i state2) Execute(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var data string
	persistence.GetDataAttribute(keyData, &data)

	i.svc.CallAPI2(data)
	return iwf.DeadEnd, nil
}

func NewState3(svc service.MyService) iwf.WorkflowState {
	return state3{svc: svc}
}

type state3 struct {
	iwf.WorkflowStateDefaults
	svc service.MyService
}

func (i state3) WaitUntil(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.AnyCommandCompletedRequest(
		iwf.NewTimerCommand("", time.Now().Add(time.Hour*24)),
		iwf.NewSignalCommand("", SignalChannelReady),
	), nil
}

func (i state3) Execute(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var data string
	persistence.GetDataAttribute(keyData, &data)
	i.svc.CallAPI3(data)

	if commandResults.Timers[0].Status == iwfidl.FIRED {
		return iwf.SingleNextState(state4{}, nil), nil
	}
	return iwf.GracefulCompletingWorkflow, nil
}

func NewState4(svc service.MyService) iwf.WorkflowState {
	return state4{svc: svc}
}

type state4 struct {
	iwf.WorkflowStateDefaultsNoWaitUntil
	svc service.MyService
}

func (i state4) Execute(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var data string
	persistence.GetDataAttribute(keyData, &data)
	i.svc.CallAPI4(data)
	return iwf.GracefulCompletingWorkflow, nil
}
