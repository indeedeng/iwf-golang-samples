This subscription workflow is to match the use case described in
* [Temporal TypeScript tutorials](https://learn.temporal.io/tutorials/typescript/subscriptions/)
* [Temporal go sample](https://github.com/temporalio/subscription-workflow-project-template-go)
* [Temporal Java Sample](https://github.com/temporalio/subscription-workflow-project-template-java)
* [Cadence Java example](https://cadenceworkflow.io/docs/concepts/workflows/#example)

#### Use case statement
Build an application for a limited time Subscription (eg a 36 month Phone plan) that satisfies these conditions:

1. When the user signs up, send a welcome email and start a free trial for **TrialPeriod**.

2. When the TrialPeriod expires, start the billing process. 
 * If the user cancels during the trial, send a trial cancellation email.

3. Billing Process:
 * As long as you have not exceeded **MaxBillingPeriods**, charge the customer for the **BillingPeriodChargeAmount**.
 * Then wait for the next **BillingPeriod**.
 * If the customer cancels during a billing period, send a subscription cancellation email.
 * If Subscription has ended normally (exceeded MaxBillingPeriods without cancellation), send a subscription ended email.

4. At any point while subscriptions are ongoing, be able to look up and change any customer's amount charged and current status and info. 

Of course, this all has to be fault tolerant, scalable to millions of customers, testable, maintainable, and observable.

#### How to run


To start a subscription workflow:
* Open http://localhost:8803/subscription/start

It will return you a **workflowId**.

The controller is hard coded to start with 20s as trial period, 10s as billing period, $100 as period charge amount for 10 max billing periods

To update the period charge amount :
* Open http://localhost:8803/subscription/updateChargeAmount?workflowId=<TheWorkflowId>&newChargeAmount=<The new amount>

To cancel the subscription:
* Open http://localhost:8803/subscription/cancel?workflowId=<TheWorkflowId>

It's recommended to use a iWF state diagram to visualize the workflow design like this:
![Subscription workflow iWF state diagram](https://user-images.githubusercontent.com/4523955/216396635-1c46df3c-e087-415a-996e-16ce47e7ccb2.png)
