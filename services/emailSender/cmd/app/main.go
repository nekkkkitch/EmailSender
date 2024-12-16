package main

import (
	"EmailSender/services/emailSender/internal/router"
	"EmailSender/services/emailSender/internal/sender"
	"EmailSender/services/emailSender/internal/service"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RouterConfig *router.Config `yaml:"router" env-prefix:"ROUTER_"`
	SenderConfig *sender.Config `yaml:"sender" env-prefix:"SENDER_"`
}

func readConfig(filename string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := readConfig("./cfg.yml")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Config file read successfully")
	log.Println(*cfg.SenderConfig)
	sender, err := sender.New(cfg.SenderConfig)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Email sender connected successfully")
	service, err := service.New(sender)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Email sender service created successfully")
	router, _ := router.New(cfg.RouterConfig, service)
	err = router.Listen()
	if err != nil {
		log.Fatalln("Failed to host router:", err.Error())
	}
	log.Printf("Router is listening...")
}
