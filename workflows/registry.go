package workflows

import (
	"github.com/indeedeng/iwf-golang-samples/workflows/basic"
	"github.com/indeedeng/iwf-golang-samples/workflows/interstate"
	"github.com/indeedeng/iwf-golang-samples/workflows/persistence"
	"github.com/indeedeng/iwf-golang-samples/workflows/signal"
	"github.com/indeedeng/iwf-golang-samples/workflows/subscription"
	"github.com/indeedeng/iwf-golang-samples/workflows/timer"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

var registry = iwf.NewRegistry()

func init() {

	svc := subscription.NewMyService()
	subscriptionWf := subscription.NewSubscriptionWorkflow(svc)

	err := registry.AddWorkflows(
		subscriptionWf,
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
