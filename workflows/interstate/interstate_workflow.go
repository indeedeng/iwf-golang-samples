package interstate

import "github.com/indeedeng/iwf-golang-sdk/iwf"

type InterStateWorkflow struct{}

const interStateChannel1 = "test-inter-state-channel-1"
const interStateChannel2 = "test-inter-state-channel-2"

func (b InterStateWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.NewStartingState(&interStateWorkflowState0{}),
		iwf.NewNonStartingState(&interStateWorkflowState1{}),
		iwf.NewNonStartingState(&interStateWorkflowState2{}),
	}
}

func (b InterStateWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return nil
}

func (b InterStateWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.NewInterstateChannelDef(interStateChannel1),
		iwf.NewInterstateChannelDef(interStateChannel2),
	}
}

func (b InterStateWorkflow) GetWorkflowType() string {
	return ""
}
