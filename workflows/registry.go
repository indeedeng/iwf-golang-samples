package workflows

import (
	"github.com/indeedeng/iwf-golang-samples/workflows/service"
	"github.com/indeedeng/iwf-golang-samples/workflows/subscription"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

var registry = iwf.NewRegistry()

func init() {

	svc := service.NewMyService()
	subscriptionWf := subscription.NewSubscriptionWorkflow(svc)

	err := registry.AddWorkflows(
		subscriptionWf,
	)
	if err != nil {
		panic(err)
	}
}

func GetRegistry() iwf.Registry {
	return registry
}
