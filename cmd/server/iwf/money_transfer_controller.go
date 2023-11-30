package iwf

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf-golang-samples/workflows/moneytransfer"
	"net/http"
	"strconv"
	"time"
)

func startMoneyTransferWorkflow(c *gin.Context) {
	fromAccount := c.Query("fromAccount")
	toAccount := c.Query("toAccount")
	amount := c.Query("amount")
	notes := c.Query("notes")

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, "must provide correct amount via URL parameter")
		return
	}

	req := moneytransfer.TransferRequest{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Notes:       notes,
		Amount:      amountInt,
	}
	wfId := fmt.Sprintf("money_transfer-%d", time.Now().Unix())

	_, err = client.StartWorkflow(c.Request.Context(), moneytransfer.MoneyTransferWorkflow{}, wfId, 3600, req, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("workflowId: %v", wfId))
	return
}
