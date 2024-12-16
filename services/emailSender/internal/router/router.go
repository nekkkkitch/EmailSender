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
	SendCode(email, first_name string) (string, error)
}

func New(cfg *Config, svc IService) (*Router, error) {
	app := fiber.New()
	router := Router{cfg: cfg, senderService: svc, app: app}
	router.app.Post("/sendemail", router.SendCode())
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
		code, err := r.senderService.SendCode(data["email"], data["name"])
		if err != nil {
			slog.Error(fmt.Sprintf("router SendCode send error: %v", err.Error()))
			c.Status(500)
			c.WriteString("Warning: message wasn't sent")
		}
		c.Status(200)
		return c.JSON(fiber.Map{"code": code})
	}
}
