package config

import (
	"login-service/helper"

	"github.com/gofiber/fiber/v2"
)

var IPPort, Net = helper.GetAddress()

var Iteung = fiber.Config{
	Prefork:       true,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "login-service",
	AppName:       "Golang Change Root",
	Network:       Net,
}
