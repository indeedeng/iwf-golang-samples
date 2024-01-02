### How to run
* Start a iWF server following the [instructions](https://github.com/indeedeng/iwf#how-to-use)
  * The easiest way is to run `docker run -p 8801:8801 -p 7233:7233 -p 8233:8233 -e AUTO_FIX_WORKER_URL=host.docker.internal --add-host host.docker.internal:host-gateway -it iworkflowio/iwf-server-lite:latest`
* Build and run this project `make bins && ./iwf-samples start`
* Start a workflow: `http://localhost:8803/polling/start?workflowId=test1&pollingCompletionThreshold=1000`
  * pollingCompletionThreshold means how many times the workflow will poll before complete the polling task C
* Signal the workflow to complete task A and B:
  * complete task A: `http://localhost:8803/polling/complete?workflowId=test1&channel=taskACompleted`
  * complete task B: `http://localhost:8803/polling/complete?workflowId=test1&channel=taskACompleted`
  * alternatively you can signal the workflow in WebUI manually 
* Watch in WebUI `http://localhost:8233/namespaces/default/workflows`
* Modify the pollingCompletionThreshold and see how the workflow complete task C automatically
