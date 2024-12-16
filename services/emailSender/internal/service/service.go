package service

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
)

type SenderService struct {
	sender IEmailSender
}

type IEmailSender interface {
	SendEmail(content []byte, reciever string) error
}

func New(sender IEmailSender) (*SenderService, error) {
	return &SenderService{sender: sender}, nil
}

func (s *SenderService) SendCode(email, first_name string) (string, error) {
	code := generateCode()
	content := []byte(fmt.Sprintf("Здравствуйте, %s. Ваш код для подтверждения: %s", first_name, code))
	err := s.sender.SendEmail(content, email)
	if err != nil {
		slog.Error(fmt.Sprintf("service SendCode error: %v", err.Error()))
		return code, err
	}
	slog.Info(fmt.Sprintf("Made code %s for user %s", code, email))
	return code, nil
}

func generateCode() string {
	n := 6
	code := ""
	alph := "1234567890"
	for range n {
		code += string([]rune(alph)[rand.IntN(len(alph))])
	}
	return code
}
