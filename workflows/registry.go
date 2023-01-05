package workflows

import (
	"github.com/indeedeng/iwf-golang-samples/workflows/basic"
	"github.com/indeedeng/iwf-golang-samples/workflows/interstate"
	"github.com/indeedeng/iwf-golang-samples/workflows/persistence"
	"github.com/indeedeng/iwf-golang-samples/workflows/signal"
	"github.com/indeedeng/iwf-golang-samples/workflows/timer"
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
