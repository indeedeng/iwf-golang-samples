package iwf

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf-golang-samples/workflows/subscription"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"net/http"
	"strconv"
)

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

func descSubscription(c *gin.Context) {
	wfId := c.Query("workflowId")
	if wfId != "" {
		wf := subscription.SubscriptionWorkflow{}
		var rpcOutput subscription.Subscription
		err := client.InvokeRPC(context.Background(), wfId, "", wf.Describe, nil, &rpcOutput)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, rpcOutput)
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
