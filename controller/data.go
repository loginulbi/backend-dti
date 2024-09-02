package controller

import (
	"login-service/config"
	"login-service/helper/at"
	"login-service/helper/atdb"
	"login-service/helper/watoken"
	"login-service/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserBio(c *fiber.Ctx, req *http.Request) error {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetLoginFromHeader(req)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		return c.Status(http.StatusForbidden).JSON(respn)
	}
	dockaryawan, err := atdb.GetOneDoc[model.Karyawan](config.Mongoconn, "karyawan", primitive.M{"nama": payload.Alias})
	if err != nil {
		dockaryawan.PhoneNumber = payload.Id
		dockaryawan.Nama = payload.Alias
		return c.Status(http.StatusOK).JSON(dockaryawan)
	}
	dockaryawan.Nama = payload.Alias
	return c.Status(http.StatusOK).JSON(dockaryawan)
}