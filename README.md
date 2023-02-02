# iwf-golang-samples

Samples for [iWF Golang SDK](https://github.com/indeedeng/iwf-golang-sdk) that runs
against [iWF server](https://github.com/indeedeng/iwf)

## Setup

1. Start a iWF server following the [instructions](https://github.com/indeedeng/iwf#how-to-run-this-server)
2. Run this project
  * To build the binary, run `make bins` 
  * To run the sample service: run `./iwf-samples start`

_Note that by default this project will listen on 8803 port(default worker port for iWF Golang SDK)_

## Product Use case samples

### Subscription workflow

This [Subscription workflow](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/subscription) is to match the use case described in
* [Temporal TypeScript tutorials](https://learn.temporal.io/tutorials/typescript/subscriptions/)
* [Temporal go sample](https://github.com/temporalio/subscription-workflow-project-template-go)
* [Temporal Java Sample](https://github.com/temporalio/subscription-workflow-project-template-java)
* [Cadence Java example](https://cadenceworkflow.io/docs/concepts/workflows/#example)

To start a subscription workflow:
* Open http://localhost:8803/subscription/start

It will return you a **workflowId**.

The controller is hard coded to start with 20s as trial period, 10s as billing period, $100 as period charge amount for 10 max billing periods

To update the period charge amount :
* Open http://localhost:8803/subscription/updateChargeAmount?workflowId=<TheWorkflowId>&newChargeAmount=<The new amount>

To cancel the subscription:
* Open http://localhost:8803/subscription/cancel?workflowId=<TheWorkflowId>

It's recommended to use a iWF state diagram to visualize the workflow design like this:

## iWF feature samples

1. [Basic IO workflow](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/basic):
   Open http://localhost:8803/basic/start in your browser. This workflow demonstrates:
    * How to start workflow with input and get output
    * How to pass input from a state to a next state
2. [Persistence workflow](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/persistence):
   Open http://localhost:8803/persistence/start in your browser. This workflow demonstrates:
    * How to use data objects to share data across workflows
    * How to use search attributes to share data and also searching for workflows
    * How to use record events API
    * How to use StateLocal to pass data from start to decide API
3. [Signal workflow](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/signal):
   Open http://localhost:8803/signal/start in your browser. This workflow demonstrates:
    * How to use signal
    * How to use AnyCommandCompleted trigger type
    * State1 start API will wait for two signals, when any of them is received, the decide API is trigger
4. [Timer workflow](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/timer):
   Open http://localhost:8803/timer/start in your browser. This workflow demonstrates:
    * How to use a durable timer
    * State1 start API will wait for a timer, when timer fires, the decide API is trigger
5. [InterstateChannel workflow](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/interstate):
   Open http://localhost:8803/interstateChannel/start in your browser. This workflow demonstrates:
    * How to use interstate channel to synchronize multi threading/in parallel workflow execution
    * State0 will go to State1 and State2
    * State1 will wait for a InterStateChannel from State2
    * State2 will send a signal and then finish as a "dead end"

Then watch the workflow in Cadence or Temporal Web UI

See [more samples in SDK integration tests](https://github.com/indeedeng/iwf-golang-sdk/tree/main/integ) for how to interact with the clients.

