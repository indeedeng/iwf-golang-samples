package subscription

import (
	"github.com/golang/mock/gomock"
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

// TODO complete the rest unit tests like in Java samples...
// https://github.com/indeedeng/iwf-java-samples/blob/main/src/test/java/io/iworkflow/workflow/subscription/SubscriptionWorkflowTest.java
