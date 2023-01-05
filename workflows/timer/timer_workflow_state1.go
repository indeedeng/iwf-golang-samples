package timer

import (
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

type timerWorkflowState1 struct{}

const TimerWorkflowState1Id = "timerWorkflowState1"

func (b timerWorkflowState1) GetStateId() string {
	return TimerWorkflowState1Id
}

func (b timerWorkflowState1) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	var i int
	err := input.Get(&i)
	if err != nil {
		return nil, err
	}
	return iwf.AllCommandsCompletedRequest(
		iwf.NewTimerCommand("", time.Now().Add(time.Duration(i)*time.Second)),
	), nil
}

func (b timerWorkflowState1) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	var i int
	err := input.Get(&i)
	if err != nil {
		return nil, err
	}
	return iwf.GracefulCompleteWorkflow(i + 1), nil
}

func (b timerWorkflowState1) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
