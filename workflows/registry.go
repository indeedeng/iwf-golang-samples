package workflows

import (
	"github.com/iworkflowio/iwf-golang-samples/workflows/basic"
	"github.com/iworkflowio/iwf-golang-samples/workflows/interstate"
	"github.com/iworkflowio/iwf-golang-samples/workflows/persistence"
	"github.com/iworkflowio/iwf-golang-samples/workflows/signal"
	"github.com/iworkflowio/iwf-golang-samples/workflows/timer"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
)

var registry = iwf.NewRegistry()

func init() {
	err := registry.AddWorkflows(
		&basic.BasicWorkflow{},
		&interstate.InterStateWorkflow{},
		&persistence.PersistenceWorkflow{},
		&signal.SignalWorkflow{},
		&timer.TimerWorkflow{},
	)
	if err != nil {
		panic(err)
	}
}

func GetRegistry() iwf.Registry {
	return registry
}
