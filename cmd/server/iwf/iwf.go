// Copyright (c) 2021 Cadence workflow OSS organization
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package iwf

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf-golang-samples/workflows"
	"github.com/indeedeng/iwf-golang-samples/workflows/basic"
	"github.com/indeedeng/iwf-golang-samples/workflows/interstate"
	"github.com/indeedeng/iwf-golang-samples/workflows/persistence"
	"github.com/indeedeng/iwf-golang-samples/workflows/signal"
	"github.com/indeedeng/iwf-golang-samples/workflows/subscription"
	"github.com/indeedeng/iwf-golang-samples/workflows/timer"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// BuildCLI is the main entry point for the iwf server
func BuildCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "iwf golang samples"
	app.Usage = "iwf golang samples"
	app.Version = "beta"

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{""},
			Usage:   "start iwf golang samples",
			Action:  start,
		},
	}
	return app
}

func start(c *cli.Context) {
	fmt.Println("start running samples")
	closeFn := startWorkflowWorker()
	// TODO improve the waiting with process signal
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	closeFn()
}

var client = iwf.NewClient(workflows.GetRegistry(), nil)
var workerService = iwf.NewWorkerService(workflows.GetRegistry(), nil)

func startWorkflowWorker() (closeFunc func()) {
	router := gin.Default()
	router.POST(iwf.WorkflowStateStartApi, apiV1WorkflowStateStart)
	router.POST(iwf.WorkflowStateDecideApi, apiV1WorkflowStateDecide)

	persInput := persistence.ExampleDataObjectModel{
		IntValue: time.Now().UnixNano(),
		StrValue: "same string for test",
		Datetime: time.Now(),
	}

	customer := subscription.Customer{
		FirstName: "Quanzheng",
		LastName:  "Long",
		Id:        "qlong",
		Email:     "qlong.seattle@gmail.com",
		Subscription: subscription.Subscription{
			TrialPeriod:         time.Second * 20,
			BillingPeriod:       time.Second * 10,
			MaxBillingPeriods:   10,
			BillingPeriodCharge: 100,
		},
	}

	router.GET("/basic/start", startWorklfow(&basic.BasicWorkflow{}, 1))
	router.GET("/interstateChannel/start", startWorklfow(&interstate.InterStateWorkflow{}, nil))
	router.GET("/persistence/start", startWorklfow(&persistence.PersistenceWorkflow{}, persInput))
	router.GET("/signal/start", startWorklfow(&signal.SignalWorkflow{}, nil))
	router.GET("/timer/start", startWorklfow(&timer.TimerWorkflow{}, 5))
	router.GET("/subscription/start", startWorklfow(&subscription.SubscriptionWorkflow{}, customer))
	router.GET("/subscription/cancel", cancelSubscription)
	router.GET("/subscription/updateChargeAmount", updateSubscriptionChargeAmount)

	wfServer := &http.Server{
		Addr:    ":" + iwf.DefaultWorkerPort,
		Handler: router,
	}
	go func() {
		if err := wfServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return func() { wfServer.Close() }
}

func startWorklfow(wf iwf.Workflow, input interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		wfId := "TestSample" + strconv.Itoa(int(time.Now().Unix()))
		runId, err := client.StartWorkflow(c.Request.Context(), wf, wfId, 3600, input, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("workflowId: %v runId: %v", wfId, runId))
		return
	}
}

func cancelSubscription(c *gin.Context) {
	wfId := c.Query("workflowId")
	if wfId != "" {
		err := client.SignalWorkflow(c.Request.Context(), &subscription.SubscriptionWorkflow{}, wfId, "", subscription.SignalCancelSubscription, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, struct{}{})
		}
		return
	}
	c.JSON(http.StatusBadRequest, "must provide workflowId via URL parameter")
}

func updateSubscriptionChargeAmount(c *gin.Context) {
	wfId := c.Query("workflowId")
	newChargeAmountStr := c.Query("newChargeAmount")
	newAmount, err := strconv.Atoi(newChargeAmountStr)

	if wfId != "" && err == nil {
		err := client.SignalWorkflow(c.Request.Context(), &subscription.SubscriptionWorkflow{}, wfId, "", subscription.SignalUpdateBillingPeriodChargeAmount, newAmount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, iwf.GetOpenApiErrorBody(err))
		} else {
			c.JSON(http.StatusOK, struct{}{})
		}
		return
	}
	c.JSON(http.StatusBadRequest, "must provide correct workflowId and newChargeAmount via URL parameter")
}

func apiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := workerService.HandleWorkflowStateStart(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}
func apiV1WorkflowStateDecide(c *gin.Context) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := workerService.HandleWorkflowStateDecide(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}
