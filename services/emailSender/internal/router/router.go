package router

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Port string `yaml:"port" env-prefix:"PORT"`
	Host string `yaml:"host" env-prefix:"HOST"`
}

type Router struct {
	cfg           *Config
	app           *fiber.App
	senderService IService
}

type IService interface {
	SendCode(email, first_name string) error
	VerifyCode(email, code string) (bool, error)
}

func New(cfg *Config, svc IService) (*Router, error) {
	app := fiber.New()
	router := Router{cfg: cfg, senderService: svc, app: app}
	router.app.Post("/sendemail", router.SendCode())
	router.app.Post("/verifycode", router.VerifyCode())
	return &router, nil
}

func (r *Router) Listen() error {
	err := r.app.Listen(r.cfg.Host + ":" + r.cfg.Port)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to listen on %v: %v", r.cfg.Host+":"+r.cfg.Port, err))
	}
	return err
}

func (r *Router) SendCode() fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := map[string]string{}
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			slog.Error(fmt.Sprintf("router SendCode unmarshal error: %v", err.Error()))
			c.Status(500)
			return err
		}
		err = r.senderService.SendCode(data["email"], data["name"])
		if err != nil {
			slog.Error(fmt.Sprintf("router SendCode send error: %v", err.Error()))
			c.Status(500)
			return err
		}
		c.Status(200)
		return nil
	}
}

func (r *Router) VerifyCode() fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := map[string]string{}
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			slog.Error(fmt.Sprintf("router VerifyCode unmarshal error: %v", err.Error()))
			c.Status(500)
			return err
		}
		isVerify, err := r.senderService.VerifyCode(data["email"], data["code"])
		if err != nil {
			slog.Error(fmt.Sprintf("router VerifyCode verify error: %v", err.Error()))
			c.Status(500)
			return err
		}
		c.Status(200)
		return c.JSON(fiber.Map{"isVerify": isVerify})
	}
}
