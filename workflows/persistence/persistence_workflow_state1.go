package persistence

import (
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

type persistenceWorkflowState1 struct {
	iwf.DefaultStateIdAndOptions
}

func (b persistenceWorkflowState1) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	var do ExampleDataObjectModel
	persistence.GetDataObject(testDataObjectKey, &do)
	if do.StrValue == "" && do.IntValue == 0 {
		input.Get(&do)
		if do.StrValue == "" || do.IntValue == 0 {
			panic("this value shouldn't be empty as we got it from start request")
		}
	} else {
		panic("this value should be empty because we haven't set it before")
	}
	persistence.SetDataObject(testDataObjectKey, do)
	persistence.SetSearchAttributeInt(testSearchAttributeInt, 1)

	return iwf.EmptyCommandRequest(), nil
}

func (b persistenceWorkflowState1) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	iv := persistence.GetSearchAttributeInt(testSearchAttributeInt)
	if iv != 1 {
		panic("this value must be 1 because it got set by Start API")
	}

	var do ExampleDataObjectModel
	persistence.GetDataObject(testDataObjectKey, &do)
	persistence.SetSearchAttributeDatetime(testSearchAttributeDatetime, do.Datetime)
	persistence.SetSearchAttributeBool(testSearchAttributeBool, true)
	return iwf.SingleNextState(persistenceWorkflowState2{}, nil), nil
}
