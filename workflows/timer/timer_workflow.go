package timer

import "github.com/indeedeng/iwf-golang-sdk/iwf"

type TimerWorkflow struct{}

func (b TimerWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.NewStartingState(&timerWorkflowState1{}),
	}
}

func (b TimerWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return nil
}

func (b TimerWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return nil
}

func (b TimerWorkflow) GetWorkflowType() string {
	return ""
}
