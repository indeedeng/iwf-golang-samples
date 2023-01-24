package signal

import "github.com/indeedeng/iwf-golang-sdk/iwf"

type SignalWorkflow struct {
	iwf.EmptyPersistenceSchema
	iwf.DefaultWorkflowType
}

const testChannelName1 = "test-channel-name-1"
const testChannelName2 = "test-channel-name-2"

func (b SignalWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(&signalWorkflowState1{}),
	}
}

func (b SignalWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return nil
}

func (b SignalWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.SignalChannelDef(testChannelName1),
		iwf.SignalChannelDef(testChannelName2),
	}
}

func (b SignalWorkflow) GetWorkflowType() string {
	return ""
}
