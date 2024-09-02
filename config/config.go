package config

import (
	"login-service/helper"
	"os"

	"github.com/gofiber/fiber/v2"
)

var PrivateKey string = os.Getenv("PRKEY")
var IPPort, Net = helper.GetAddress()

var Iteung = fiber.Config{
	Prefork:       true,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "login-service",
	AppName:       "Golang Change Root",
	Network:       Net,
}
