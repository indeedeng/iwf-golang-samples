package iwf

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf-golang-samples/workflows/microservices"
	"net/http"
)

func startMicroserviceWorkflow(c *gin.Context) {
	wfId := c.Query("workflowId")
	if wfId != "" {
		wf := microservices.OrchestrationWorkflow{}
		runId, err := client.StartWorkflow(c.Request.Context(), wf, wfId, 3600, "test initial data", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("workflowId: %v runId: %v", wfId, runId))
		return
	}
	c.JSON(http.StatusBadRequest, "must provide workflowId via URL parameter")
}

func signalMicroserviceWorkflow(c *gin.Context) {
	wfId := c.Query("workflowId")
	if wfId != "" {
		wf := microservices.OrchestrationWorkflow{}
		err := client.SignalWorkflow(context.Background(), wf, wfId, "", microservices.SignalChannelReady, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, struct{}{})
		}
		return
	}
	c.JSON(http.StatusBadRequest, "must provide workflowId via URL parameter")
}

func swapDataMicroserviceWorkflow(c *gin.Context) {
	wfId := c.Query("workflowId")
	newData := c.Query("data")
	if wfId != "" {
		wf := microservices.OrchestrationWorkflow{}
		var output string
		err := client.InvokeRPC(context.Background(), wfId, "", wf.Swap, newData, &output)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, output)
		}
		return
	}
	c.JSON(http.StatusBadRequest, "must provide workflowId via URL parameter")
}
