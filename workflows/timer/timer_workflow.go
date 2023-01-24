package timer

import "github.com/indeedeng/iwf-golang-sdk/iwf"

type TimerWorkflow struct {
	iwf.EmptyCommunicationSchema
	iwf.EmptyPersistenceSchema
	iwf.DefaultWorkflowType
}

func (b TimerWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(&timerWorkflowState1{}),
	}
}
