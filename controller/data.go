package controller

import (
	"login-service/config"
	"login-service/helper/at"
	"login-service/helper/atdb"
	"login-service/helper/watoken"
	"login-service/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserBio(ctx *fiber.Ctx) error {
	// Ambil login header menggunakan GetLoginFromHeader
	var dbname = "hris"
	loginSecret, err := at.GetLoginFromHeader(ctx)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = err.Error()
		respn.Location = "Missing or invalid login header"
		return ctx.Status(fiber.StatusForbidden).JSON(respn)
	}

	// Decode token menggunakan loginSecret yang didapat dari header
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, loginSecret)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = loginSecret
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		return ctx.Status(fiber.StatusForbidden).JSON(respn)
	}

	// Ambil data karyawan dari database
	dockaryawan, err := atdb.GetOneDoc[model.Karyawan](config.DBPresensi(dbname), "karyawan", bson.M{"nama": payload.Alias})
	if err != nil {
		// Jika tidak ditemukan, buat data karyawan baru berdasarkan payload
		dockaryawan.PhoneNumber = payload.Id
		dockaryawan.Nama = payload.Alias
		return ctx.Status(fiber.StatusOK).JSON(dockaryawan)
	}

	// Jika ditemukan, kembalikan data karyawan
	dockaryawan.Nama = payload.Alias
	return ctx.Status(fiber.StatusOK).JSON(dockaryawan)
}
