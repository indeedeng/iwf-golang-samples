package workflows

import (
	"github.com/indeedeng/iwf-golang-samples/workflows/engagement"
	"github.com/indeedeng/iwf-golang-samples/workflows/microservices"
	"github.com/indeedeng/iwf-golang-samples/workflows/moneytransfer"
	"github.com/indeedeng/iwf-golang-samples/workflows/service"
	"github.com/indeedeng/iwf-golang-samples/workflows/subscription"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
)

var registry = iwf.NewRegistry()

func init() {

	svc := service.NewMyService()

	err := registry.AddWorkflows(
		subscription.NewSubscriptionWorkflow(svc),
		engagement.NewEngagementWorkflow(svc),
		microservices.NewMicroserviceOrchestrationWorkflow(svc),
		moneytransfer.NewMoneyTransferWorkflow(svc),
	)
	if err != nil {
		panic(err)
	}
}

func GetRegistry() iwf.Registry {
	return registry
}
