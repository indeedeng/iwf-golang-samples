package subscription

import (
	"github.com/golang/mock/gomock"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"github.com/indeedeng/iwf-golang-sdk/iwftest"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// mockgen -source=workflows/subscription/my_service.go -destination=workflows/subscription/my_service_mock.go --package=subscription

var testCustomer = Customer{
	FirstName: "Quanzheng",
	LastName:  "Long",
	Id:        "123",
	Email:     "qlong.seattle@gmail.com",
	Subscription: Subscription{
		BillingPeriod:       time.Second,
		MaxBillingPeriods:   10,
		TrialPeriod:         time.Second * 2,
		BillingPeriodCharge: 100,
	},
}

var testCustomerObj = iwftest.NewTestObject(testCustomer)

var mockWfCtx *iwftest.MockWorkflowContext
var mockPersistence *iwftest.MockPersistence
var mockCommunication *iwftest.MockCommunication
var emptyCmdResults = iwf.CommandResults{}
var emptyObj = iwftest.NewTestObject(nil)
var mockSvc *MockMyService

func beforeEach(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockSvc = NewMockMyService(ctrl)
	mockWfCtx = iwftest.NewMockWorkflowContext(ctrl)
	mockPersistence = iwftest.NewMockPersistence(ctrl)
	mockCommunication = iwftest.NewMockCommunication(ctrl)
}

func TestInitState_WaitUntil(t *testing.T) {
	beforeEach(t)

	state := NewInitState()

	mockPersistence.EXPECT().SetDataAttribute(keyCustomer, testCustomer)
	cmdReq, err := state.WaitUntil(mockWfCtx, testCustomerObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.EmptyCommandRequest(), cmdReq)
}

func TestInitState_Execute(t *testing.T) {
	beforeEach(t)

	state := NewInitState()
	input := iwftest.NewTestObject(testCustomer)

	decision, err := state.Execute(mockWfCtx, input, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.MultiNextStates(
		trialState{}, cancelState{}, updateChargeAmountState{},
	), decision)
}

func TestTrialState_WaitUntil(t *testing.T) {
	beforeEach(t)

	state := NewTrialState(mockSvc)

	mockSvc.EXPECT().sendEmail(testCustomer.Email, gomock.Any(), gomock.Any())
	mockPersistence.EXPECT().GetDataAttribute(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	cmdReq, err := state.WaitUntil(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	firingTime := cmdReq.Commands[0].TimerCommand.FiringUnixTimestampSeconds
	assert.Equal(t, iwf.AllCommandsCompletedRequest(
		iwf.NewTimerCommand("", time.Unix(firingTime, 0)),
	), cmdReq)
}

func TestTrialState_Execute(t *testing.T) {
	beforeEach(t)

	state := NewTrialState(mockSvc)

	mockPersistence.EXPECT().SetDataAttribute(keyBillingPeriodNum, 0)

	decision, err := state.Execute(mockWfCtx, emptyObj, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.SingleNextState(
		chargeCurrentBillState{}, nil,
	), decision)
}

func TestChargeCurrentBillStateStart_waitForDuration(t *testing.T) {
	beforeEach(t)

	state := NewChargeCurrentBillState(mockSvc)

	mockPersistence.EXPECT().GetDataAttribute(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().GetDataAttribute(keyBillingPeriodNum, gomock.Any()).SetArg(1, 0)
	mockPersistence.EXPECT().SetDataAttribute(keyBillingPeriodNum, 1)

	cmdReq, err := state.WaitUntil(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	cmd := cmdReq.Commands[0]
	assert.Equal(t, iwf.AllCommandsCompletedRequest(iwf.NewTimerCommand("", time.Unix(cmd.TimerCommand.FiringUnixTimestampSeconds, 0))), cmdReq)
}

func TestChargeCurrentBillStateStart_subscriptionOver(t *testing.T) {
	beforeEach(t)

	state := NewChargeCurrentBillState(mockSvc)

	mockPersistence.EXPECT().GetDataAttribute(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().GetDataAttribute(keyBillingPeriodNum, gomock.Any()).SetArg(1, testCustomer.Subscription.MaxBillingPeriods)
	mockPersistence.EXPECT().SetStateExecutionLocal(subscriptionOverKey, true)

	cmdReq, err := state.WaitUntil(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.EmptyCommandRequest(), cmdReq)
}

func TestChargeCurrentBillStateDecide_subscriptionNotOver(t *testing.T) {
	beforeEach(t)

	state := NewChargeCurrentBillState(mockSvc)

	mockPersistence.EXPECT().GetDataAttribute(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().GetStateExecutionLocal(subscriptionOverKey, gomock.Any())
	mockSvc.EXPECT().chargeUser(testCustomer.Email, testCustomer.Id, testCustomer.Subscription.BillingPeriodCharge)

	decision, err := state.Execute(mockWfCtx, emptyObj, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.SingleNextState(&chargeCurrentBillState{}, nil), decision)
}

func TestChargeCurrentBillStateDecide_subscriptionOver(t *testing.T) {
	beforeEach(t)

	state := NewChargeCurrentBillState(mockSvc)

	mockPersistence.EXPECT().GetDataAttribute(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().GetStateExecutionLocal(subscriptionOverKey, gomock.Any()).SetArg(1, true)
	mockSvc.EXPECT().sendEmail(testCustomer.Email, gomock.Any(), gomock.Any())

	decision, err := state.Execute(mockWfCtx, emptyObj, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.ForceCompletingWorkflow, decision)
}

func TestUpdateChargeAmountState_WaitUntil(t *testing.T) {
	beforeEach(t)

	state := NewUpdateChargeAmountState()

	cmdReq, err := state.WaitUntil(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.AllCommandsCompletedRequest(iwf.NewSignalCommand("", SignalUpdateBillingPeriodChargeAmount)), cmdReq)
}

func TestUpdateChargeAmountState_Execute(t *testing.T) {
	beforeEach(t)

	state := NewUpdateChargeAmountState()

	cmdResults := iwf.CommandResults{
		Signals: []iwf.SignalCommandResult{
			{
				ChannelName: SignalUpdateBillingPeriodChargeAmount,
				SignalValue: iwftest.NewTestObject(200),
				Status:      iwfidl.RECEIVED,
			},
		},
	}

	updatedCustomer := testCustomer
	updatedCustomer.Subscription.BillingPeriodCharge = 200

	mockPersistence.EXPECT().GetDataAttribute(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().SetDataAttribute(keyCustomer, updatedCustomer)

	decision, err := state.Execute(mockWfCtx, emptyObj, cmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.SingleNextState(&updateChargeAmountState{}, nil), decision)
}

func TestCancelState_WaitUntil(t *testing.T) {
	beforeEach(t)

	state := NewCancelState(mockSvc)

	cmdReq, err := state.WaitUntil(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.AllCommandsCompletedRequest(iwf.NewSignalCommand("", SignalCancelSubscription)), cmdReq)
}

func TestCancelState_Execute(t *testing.T) {
	beforeEach(t)

	state := NewCancelState(mockSvc)

	mockPersistence.EXPECT().GetDataAttribute(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockSvc.EXPECT().sendEmail(testCustomer.Email, gomock.Any(), gomock.Any())

	decision, err := state.Execute(mockWfCtx, emptyObj, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.ForceCompletingWorkflow, decision)
}
