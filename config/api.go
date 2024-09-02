package config

import "os"

var WAKeyword string = os.Getenv("WAQRKEYWORD")

var WebhookURL string = os.Getenv("WEBHOOKURL")

var WebhookSecret string = os.Getenv("WEBHOOKSECRET")

var WAPhoneNumber string = os.Getenv("WAPHONENUMBER")

var WAAPIQRLogin string = "https://api.wa.my.id/api/whatsauth/request"

var WAAPIMessage string = "https://api.wa.my.id/api/send/message/text"

var WAAPIDocMessage string = "https://api.wa.my.id/api/send/message/document"

var WAAPIGetToken string = "https://api.wa.my.id/api/signup"

var WAAPIGetDevice string = "https://api.wa.my.id/api/device/"

var PublicKeyWhatsAuth string 

var WAAPIToken string