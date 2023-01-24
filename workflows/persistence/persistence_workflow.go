package persistence

import (
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type PersistenceWorkflow struct {
	iwf.EmptyCommunicationSchema
	iwf.DefaultWorkflowType
}

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
		iwf.StartingStateDef(&persistenceWorkflowState1{}),
		iwf.NonStartingStateDef(&persistenceWorkflowState2{}),
	}
}

func (b PersistenceWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.DataObjectDef(testDataObjectKey),
		iwf.SearchAttributeDef(testSearchAttributeInt, iwfidl.INT),
		iwf.SearchAttributeDef(testSearchAttributeDatetime, iwfidl.DATETIME),
		iwf.SearchAttributeDef(testSearchAttributeBool, iwfidl.BOOL),
		iwf.SearchAttributeDef(testSearchAttributeDouble, iwfidl.DOUBLE),
		iwf.SearchAttributeDef(testSearchAttributeText, iwfidl.TEXT),
		iwf.SearchAttributeDef(testSearchAttributeKeyword, iwfidl.KEYWORD),
	}
}
