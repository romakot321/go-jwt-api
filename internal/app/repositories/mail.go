package repositories

import (
  "log"
)

type MailRepository interface {
  SendIPChangedWarning(email, oldIP, newIP string)
}

type mailRepository struct {

}

func (s mailRepository) SendIPChangedWarning(email string, oldIP string, newIP string) {
  log.Print("Sending ip change warning...")
}

func NewMailRepository() MailRepository {
  return &mailRepository{}
}
