package subscription

import "fmt"

type MyService interface {
	sendEmail(recipient, subject, content string)
	chargeUser(email, customerId string, amount int)
}

type myServiceImpl struct{}

func (m myServiceImpl) sendEmail(recipient, subject, content string) {
	fmt.Printf("sending an email to %v, title: %v, content: %v \n", recipient, subject, content)
}

func (m myServiceImpl) chargeUser(email, customerId string, amount int) {
	fmt.Printf("charege user customerId[%v] email[%v] for $%v \n", customerId, email, amount)
}

func NewMyService() MyService {
	return &myServiceImpl{}
}
