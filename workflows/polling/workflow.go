package polling

import (
	"github.com/indeedeng/iwf-golang-samples/workflows/service"
	"github.com/indeedeng/iwf-golang-sdk/gen/iwfidl"
	"github.com/indeedeng/iwf-golang-sdk/iwf"
	"time"
)

func NewPollingWorkflow(svc service.MyService) iwf.ObjectWorkflow {

	return &PollingWorkflow{
		svc: svc,
	}
}

const (
	dataAttrTaskACompleted = "taskACompleted"
	dataAttrTaskBCompleted = "taskBCompleted"
	dataAttrTaskCCompleted = "taskCCompleted"

	dataAttrCurrPolls = "currPolls" // tracks how many polls have been done

	SignalChannelTaskACompleted = "taskACompleted"
	SignalChannelTaskBCompleted = "taskBCompleted"
)

type PollingWorkflow struct {
	iwf.WorkflowDefaults

	svc service.MyService
}

func (e PollingWorkflow) GetWorkflowStates() []iwf.StateDef {
	return []iwf.StateDef{
		iwf.StartingStateDef(&checkAndCompleteState{svc: e.svc}),
	}
}

func (e PollingWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return []iwf.PersistenceFieldDef{
		iwf.DataAttributeDef(dataAttrTaskACompleted),
		iwf.DataAttributeDef(dataAttrTaskBCompleted),
		iwf.DataAttributeDef(dataAttrTaskCCompleted),
		iwf.DataAttributeDef(dataAttrCurrPolls),
	}
}

func (e PollingWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.SignalChannelDef(SignalChannelTaskACompleted),
		iwf.SignalChannelDef(SignalChannelTaskBCompleted),
	}
}

type checkAndCompleteState struct {
	iwf.WorkflowStateDefaults
	svc service.MyService
}

func (i checkAndCompleteState) WaitUntil(
	ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication,
) (*iwf.CommandRequest, error) {
	var taskACompleted bool
	persistence.GetDataAttribute(dataAttrTaskACompleted, &taskACompleted)
	var taskBCompleted bool
	persistence.GetDataAttribute(dataAttrTaskBCompleted, &taskBCompleted)
	var taskCCompleted bool
	persistence.GetDataAttribute(dataAttrTaskCCompleted, &taskCCompleted)

	var commands []iwf.Command
	if !taskACompleted {
		commands = append(commands, iwf.NewSignalCommand("", SignalChannelTaskACompleted))
	}
	if !taskBCompleted {
		commands = append(commands, iwf.NewSignalCommand("", SignalChannelTaskBCompleted))
	}

	if !taskCCompleted {
		commands = append(commands, iwf.NewTimerCommand("", time.Now().Add(time.Second*2)))
	}

	return iwf.AnyCommandCompletedRequest(commands...), nil
}

func (i checkAndCompleteState) Execute(
	ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence,
	communication iwf.Communication,
) (*iwf.StateDecision, error) {
	var taskACompleted bool
	persistence.GetDataAttribute(dataAttrTaskACompleted, &taskACompleted)
	var taskBCompleted bool
	persistence.GetDataAttribute(dataAttrTaskBCompleted, &taskBCompleted)
	var taskCCompleted bool
	persistence.GetDataAttribute(dataAttrTaskCCompleted, &taskCCompleted)

	var maxPollsRequired int
	input.Get(&maxPollsRequired)

	if !taskCCompleted {
		i.svc.CallAPI1("calling API1 for polling service C")

		var currPolls int
		persistence.GetDataAttribute(dataAttrCurrPolls, &currPolls)
		if currPolls >= maxPollsRequired {
			taskCCompleted = true
			persistence.SetDataAttribute(dataAttrTaskCCompleted, true)
		}
		persistence.SetDataAttribute(dataAttrCurrPolls, currPolls+1)
	}

	for _, signal := range commandResults.Signals {
		switch signal.ChannelName {
		case SignalChannelTaskACompleted:
			if signal.Status == iwfidl.RECEIVED {
				taskACompleted = true
				persistence.SetDataAttribute(dataAttrTaskACompleted, true)
			}
		case SignalChannelTaskBCompleted:
			if signal.Status == iwfidl.RECEIVED {
				taskBCompleted = true
				persistence.SetDataAttribute(dataAttrTaskBCompleted, true)
			}
		}
	}

	if taskACompleted && taskBCompleted && taskCCompleted {
		return iwf.GracefulCompletingWorkflow, nil
	}

	// loop back to check
	return iwf.SingleNextState(checkAndCompleteState{}, maxPollsRequired), nil
}
