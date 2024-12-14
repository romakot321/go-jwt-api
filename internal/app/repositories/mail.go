package repositories

type MailRepository interface {
  SendIPChangedWarning(email, oldIP, newIP string)
}

type mailRepository struct {

}

func (s mailRepository) SendIPChangedWarning(email string, oldIP string, newIP string) {

}

func NewMailRepository() MailRepository {
  return &mailRepository{}
}
