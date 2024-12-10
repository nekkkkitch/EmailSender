package service

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
)

type SenderService struct {
	sender IEmailSender
	redis  IRedis
}

type IEmailSender interface {
	SendEmail(content []byte, reciever string) error
}

type IRedis interface {
	GetCode(email string) (string, error)
	PutCode(email, code string) error
}

func New(sender IEmailSender, redis IRedis) (*SenderService, error) {
	return &SenderService{sender: sender, redis: redis}, nil
}

func (s *SenderService) SendCode(email, first_name string) error {
	code := generateCode()
	content := []byte(fmt.Sprintf("Здравствуйте, %s. Ваш код для подтверждения: %s", first_name, code))
	err := s.sender.SendEmail(content, email)
	if err != nil {
		slog.Error(fmt.Sprintf("service SendCode error: %v", err.Error()))
		return err
	}
	slog.Info(fmt.Sprintf("Made code %s for user %s", code, email))
	err = s.redis.PutCode(email, code)
	if err != nil {
		slog.Error(fmt.Sprintf("service SendCode error: %v", err.Error()))
		return err
	}
	return nil
}

func (s *SenderService) VerifyCode(email, code string) (bool, error) {
	dbCode, err := s.redis.GetCode(email)
	if err != nil {
		slog.Error(fmt.Sprintf("service VerifyCode error: %v", err.Error()))
		return false, err
	}
	return code == dbCode, nil
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
