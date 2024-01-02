package iwf

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf-golang-samples/workflows/polling"
	"net/http"
	"strconv"
)

func startPollingWorkflow(c *gin.Context) {
	wfId := c.Query("workflowId")
	pollingCompletionThreshold := c.Query("pollingCompletionThreshold")

	pollingCompletionThresholdInt, err := strconv.Atoi(pollingCompletionThreshold)
	if err != nil {
		c.JSON(http.StatusBadRequest, "must provide correct pollingCompletionThreshold via URL parameter")
		return
	}

	_, err = client.StartWorkflow(c.Request.Context(), polling.PollingWorkflow{}, wfId, 0, pollingCompletionThresholdInt, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("workflowId: %v is started", wfId))
	return
}

func signalPollingWorkflow(c *gin.Context) {
	wfId := c.Query("workflowId")
	channel := c.Query("channel")

	err := client.SignalWorkflow(c.Request.Context(), polling.PollingWorkflow{}, wfId, "", channel, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("workflowId: %v is signal", wfId))
	return
}
