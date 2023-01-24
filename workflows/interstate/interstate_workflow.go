package interstate

import "github.com/indeedeng/iwf-golang-sdk/iwf"

type InterStateWorkflow struct {
	iwf.EmptyPersistenceSchema
	iwf.DefaultWorkflowType
}

const interStateChannel1 = "test-inter-state-channel-1"
const interStateChannel2 = "test-inter-state-channel-2"

func (b InterStateWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(&interStateWorkflowState0{}),
		iwf.NonStartingStateDef(&interStateWorkflowState1{}),
		iwf.NonStartingStateDef(&interStateWorkflowState2{}),
	}
}

func (b InterStateWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.InterstateChannelDef(interStateChannel1),
		iwf.InterstateChannelDef(interStateChannel2),
	}
}
