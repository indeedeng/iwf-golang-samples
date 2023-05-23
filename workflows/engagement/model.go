package engagement

type Status string

const (
	StatusInitiated Status = "Initiated"
	StatusAccepted  Status = "Accepted"
	StatusDeclined  Status = "Declined"
)

type EngagementInput struct {
	EmployerId  string
	JobSeekerId string
	Notes       string
}

type EngagementDescription struct {
	EmployerId    string
	JobSeekerId   string
	Notes         string
	CurrentStatus Status
}
