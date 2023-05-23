package service

import "fmt"

type MyService interface {
	SendEmail(recipient, subject, content string)
	ChargeUser(email, customerId string, amount int)
	UpdateExternalSystem(message string)
}

type myServiceImpl struct{}

func (m myServiceImpl) UpdateExternalSystem(message string) {
	fmt.Println("Update external system(like via RPC, or sending Kafka message or database):", message)
}

func (m myServiceImpl) SendEmail(recipient, subject, content string) {
	fmt.Printf("sending an email to %v, title: %v, content: %v \n", recipient, subject, content)
}

func (m myServiceImpl) ChargeUser(email, customerId string, amount int) {
	fmt.Printf("charege user customerId[%v] email[%v] for $%v \n", customerId, email, amount)
}

func NewMyService() MyService {
	return &myServiceImpl{}
}
