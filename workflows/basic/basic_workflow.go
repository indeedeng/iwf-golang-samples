package basic

import "github.com/indeedeng/iwf-golang-sdk/iwf"

type BasicWorkflow struct{}

func (b BasicWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.NewStartingState(&basicWorkflowState1{}),
		iwf.NewNonStartingState(&basicWorkflowState2{}),
	}
}

func (b BasicWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return nil
}

func (b BasicWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return nil
}

func (b BasicWorkflow) GetWorkflowType() string {
	return ""
}
