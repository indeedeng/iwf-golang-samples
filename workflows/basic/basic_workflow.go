package basic

import "github.com/indeedeng/iwf-golang-sdk/iwf"

type BasicWorkflow struct {
	iwf.EmptyCommunicationSchema
	iwf.EmptyPersistenceSchema
	iwf.DefaultWorkflowType
}

func (b BasicWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(&basicWorkflowState1{}),
		iwf.NonStartingStateDef(&basicWorkflowState2{}),
	}
}
