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

### [JobSeeker Engagement workflow](./workflows/engagement)
<img width="709" alt="Screenshot 2023-04-21 at 8 53 25 AM" src="https://user-images.githubusercontent.com/4523955/233680837-6a6267a0-4b31-419e-87f0-667bb48582d1.png">
This engagement workflow is for: 

* An engagement is initiated by an employer to reach out to a jobSeeker(via email/SMS/etc)
* The jobSeeker could respond with decline or accept
* If jobSeeker doesn't respond, it will get reminder
* An engagement can change from declined to accepted, but cannot change from accepted to declined


### [Subscription](./workflows/subscription) workflow

This [Subscription workflow](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/subscription) (with unit tests) is to match the use case described in
* [Temporal TypeScript tutorials](https://learn.temporal.io/tutorials/typescript/subscriptions/)
* [Temporal go sample](https://github.com/temporalio/subscription-workflow-project-template-go)
* [Temporal Java Sample](https://github.com/temporalio/subscription-workflow-project-template-java)
* [Cadence Java example](https://cadenceworkflow.io/docs/concepts/workflows/#example)

