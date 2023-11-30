### How to run
* start a iWF server following the [instructions](https://github.com/indeedeng/iwf#how-to-use)
* build and run this project `make bins && ./iwf-samples start`
* start a workflow: `http://localhost:8803/moneytransfer/start?fromAccount=test1&toAccount=test2&amount=100&notes=hello`
* watch in WebUI `http://localhost:8233/namespaces/default/workflows`
* modify the workflow code to try injecting some errors, and shorten the retry, to see what will happen
