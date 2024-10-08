package url

import (
	"login-service/controller"

	"github.com/gofiber/fiber/v2"
)

func Web(page *fiber.App) {
	// page.Get("/", controller.Homepage)
	// page.Get("/ip", controller.GetIPServer)
	// page.Get("/whatsauth/refreshtoken", controller.RefreshWAToken)

	// page.Post("/whatsauth/webhook", controller.WhatsAuthReceiver)

	// page.Get("/auth/phonenumber/:login", controller.GetPhoneNumber)

	page.Post("/auth/users", controller.AuthUser)
	page.Get("/data/karyawan", func(c *fiber.Ctx) error {
		return controller.GetUserBio(c)
	})
	page.Get("/data/tes", func(c *fiber.Ctx) error {
		return controller.GetAllDataKaryawan(c)
	})
}
