package persistence

import (
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type PersistenceWorkflow struct{}

const (
	testDataObjectKey = "test-data-object"

	testSearchAttributeInt      = "CustomIntField"
	testSearchAttributeDatetime = "CustomDatetimeField"
	testSearchAttributeBool     = "CustomBoolField"
	testSearchAttributeDouble   = "CustomDoubleField"
	testSearchAttributeText     = "CustomStringField"
	testSearchAttributeKeyword  = "CustomKeywordField"
)

func (b PersistenceWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.NewStartingState(&persistenceWorkflowState1{}),
		iwf.NewNonStartingState(&persistenceWorkflowState2{}),
	}
}

func (b PersistenceWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.NewDataObjectDef(testDataObjectKey),
		iwf.NewSearchAttributeDef(testSearchAttributeInt, iwfidl.INT),
		iwf.NewSearchAttributeDef(testSearchAttributeDatetime, iwfidl.DATETIME),
		iwf.NewSearchAttributeDef(testSearchAttributeBool, iwfidl.BOOL),
		iwf.NewSearchAttributeDef(testSearchAttributeDouble, iwfidl.DOUBLE),
		iwf.NewSearchAttributeDef(testSearchAttributeText, iwfidl.TEXT),
		iwf.NewSearchAttributeDef(testSearchAttributeKeyword, iwfidl.KEYWORD),
	}
}

func (b PersistenceWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return nil
}

func (b PersistenceWorkflow) GetWorkflowType() string {
	return ""
}
