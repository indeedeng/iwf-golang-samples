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

func TestInitState_Start(t *testing.T) {
	beforeEach(t)

	state := NewInitState()

	mockPersistence.EXPECT().SetDataObject(keyCustomer, testCustomer)
	cmdReq, err := state.Start(mockWfCtx, testCustomerObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.EmptyCommandRequest(), cmdReq)
}

func TestInitState_Decide(t *testing.T) {
	beforeEach(t)

	state := NewInitState()
	input := iwftest.NewTestObject(testCustomer)

	decision, err := state.Decide(mockWfCtx, input, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.MultiNextStates(
		trialState{}, cancelState{}, updateChargeAmountState{},
	), decision)
}

func TestTrialState_Start(t *testing.T) {
	beforeEach(t)

	state := NewTrialState(mockSvc)

	mockSvc.EXPECT().sendEmail(testCustomer.Email, gomock.Any(), gomock.Any())
	mockPersistence.EXPECT().GetDataObject(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	cmdReq, err := state.Start(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	firingTime := cmdReq.Commands[0].TimerCommand.FiringUnixTimestampSeconds
	assert.Equal(t, iwf.AllCommandsCompletedRequest(
		iwf.NewTimerCommand("", time.Unix(firingTime, 0)),
	), cmdReq)
}

func TestTrialState_Decide(t *testing.T) {
	beforeEach(t)

	state := NewTrialState(mockSvc)

	mockPersistence.EXPECT().SetDataObject(keyBillingPeriodNum, 0)

	decision, err := state.Decide(mockWfCtx, emptyObj, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.SingleNextState(
		chargeCurrentBillState{}, nil,
	), decision)
}

func TestChargeCurrentBillStateStart_waitForDuration(t *testing.T) {
	beforeEach(t)

	state := NewChargeCurrentBillState(mockSvc)

	mockPersistence.EXPECT().GetDataObject(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().GetDataObject(keyBillingPeriodNum, gomock.Any()).SetArg(1, 0)
	mockPersistence.EXPECT().SetDataObject(keyBillingPeriodNum, 1)

	cmdReq, err := state.Start(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	cmd := cmdReq.Commands[0]
	assert.Equal(t, iwf.AllCommandsCompletedRequest(iwf.NewTimerCommand("", time.Unix(cmd.TimerCommand.FiringUnixTimestampSeconds, 0))), cmdReq)
}

func TestChargeCurrentBillStateStart_subscriptionOver(t *testing.T) {
	beforeEach(t)

	state := NewChargeCurrentBillState(mockSvc)

	mockPersistence.EXPECT().GetDataObject(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().GetDataObject(keyBillingPeriodNum, gomock.Any()).SetArg(1, testCustomer.Subscription.MaxBillingPeriods)
	mockPersistence.EXPECT().SetStateLocal(subscriptionOverKey, true)

	cmdReq, err := state.Start(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.EmptyCommandRequest(), cmdReq)
}

func TestChargeCurrentBillStateDecide_subscriptionNotOver(t *testing.T) {
	beforeEach(t)

	state := NewChargeCurrentBillState(mockSvc)

	mockPersistence.EXPECT().GetDataObject(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().GetStateLocal(subscriptionOverKey, gomock.Any())
	mockSvc.EXPECT().chargeUser(testCustomer.Email, testCustomer.Id, testCustomer.Subscription.BillingPeriodCharge)

	decision, err := state.Decide(mockWfCtx, emptyObj, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.SingleNextState(&chargeCurrentBillState{}, nil), decision)
}

func TestChargeCurrentBillStateDecide_subscriptionOver(t *testing.T) {
	beforeEach(t)

	state := NewChargeCurrentBillState(mockSvc)

	mockPersistence.EXPECT().GetDataObject(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().GetStateLocal(subscriptionOverKey, gomock.Any()).SetArg(1, true)
	mockSvc.EXPECT().sendEmail(testCustomer.Email, gomock.Any(), gomock.Any())

	decision, err := state.Decide(mockWfCtx, emptyObj, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.ForceCompletingWorkflow, decision)
}

func TestUpdateChargeAmountState_Start(t *testing.T) {
	beforeEach(t)

	state := NewUpdateChargeAmountState()

	cmdReq, err := state.Start(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.AllCommandsCompletedRequest(iwf.NewSignalCommand("", SignalUpdateBillingPeriodChargeAmount)), cmdReq)
}

func TestUpdateChargeAmountState_Decide(t *testing.T) {
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

	mockPersistence.EXPECT().GetDataObject(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockPersistence.EXPECT().SetDataObject(keyCustomer, updatedCustomer)

	decision, err := state.Decide(mockWfCtx, emptyObj, cmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.SingleNextState(&updateChargeAmountState{}, nil), decision)
}

func TestCancelState_Start(t *testing.T) {
	beforeEach(t)

	state := NewCancelState(mockSvc)

	cmdReq, err := state.Start(mockWfCtx, emptyObj, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.AllCommandsCompletedRequest(iwf.NewSignalCommand("", SignalCancelSubscription)), cmdReq)
}

func TestCancelState_Decide(t *testing.T) {
	beforeEach(t)

	state := NewCancelState(mockSvc)

	mockPersistence.EXPECT().GetDataObject(keyCustomer, gomock.Any()).SetArg(1, testCustomer)
	mockSvc.EXPECT().sendEmail(testCustomer.Email, gomock.Any(), gomock.Any())

	decision, err := state.Decide(mockWfCtx, emptyObj, emptyCmdResults, mockPersistence, mockCommunication)
	assert.Nil(t, err)
	assert.Equal(t, iwf.ForceCompletingWorkflow, decision)
}
