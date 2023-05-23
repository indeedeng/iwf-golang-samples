package iwf

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf-golang-samples/workflows/engagement"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"net/http"
	"strings"
)

func descEngagement(c *gin.Context) {
	wfId := c.Query("workflowId")
	if wfId != "" {
		wf := engagement.EngagementWorkflow{}
		var rpcOutput engagement.EngagementDescription
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

func optOutReminder(c *gin.Context) {
	wfId := c.Query("workflowId")
	if wfId != "" {
		wf := engagement.EngagementWorkflow{}
		err := client.SignalWorkflow(context.Background(), wf, wfId, "", engagement.SignalChannelOptOutReminder, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, struct{}{})
		}
		return
	}
	c.JSON(http.StatusBadRequest, "must provide workflowId via URL parameter")
}

func declineEngagement(c *gin.Context) {
	wfId := c.Query("workflowId")
	if wfId != "" {
		wf := engagement.EngagementWorkflow{}
		err := client.InvokeRPC(context.Background(), wfId, "", wf.Decline, nil, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, struct{}{})
		}
		return
	}
	c.JSON(http.StatusBadRequest, "must provide workflowId via URL parameter")
}

func acceptEngagement(c *gin.Context) {
	wfId := c.Query("workflowId")
	if wfId != "" {
		wf := engagement.EngagementWorkflow{}
		err := client.InvokeRPC(context.Background(), wfId, "", wf.Accept, nil, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, struct{}{})
		}
		return
	}
	c.JSON(http.StatusBadRequest, "must provide workflowId via URL parameter")
}

func listEngagements(c *gin.Context) {
	query := c.Query("query")
	if query != "" {
		if strings.HasPrefix(query, "'") {
			query = strings.Trim(query, "'")
		}
		resp, err := client.SearchWorkflow(context.Background(), iwfidl.WorkflowSearchRequest{
			Query: query,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, resp)
		}
		return
	}
	c.JSON(http.StatusBadRequest, "must provide workflowId via URL parameter")
}
