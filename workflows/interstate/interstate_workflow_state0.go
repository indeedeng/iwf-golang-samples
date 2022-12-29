package interstate

import (
	"github.com/iworkflowio/iwf-golang-sdk/gen/iwfidl"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
)

type interStateWorkflowState0 struct{}

const InterStateWorkflowState0Id = "interStateWorkflowState0"

func (b interStateWorkflowState0) GetStateId() string {
	return InterStateWorkflowState0Id
}

func (b interStateWorkflowState0) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return iwf.EmptyCommandRequest(), nil
}

func (b interStateWorkflowState0) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	return iwf.MultiNextStatesByStateIds(interStateWorkflowState1Id, interStateWorkflowState2Id), nil
}

func (b interStateWorkflowState0) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
