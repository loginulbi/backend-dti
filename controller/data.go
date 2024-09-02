package controller

import (
	"login-service/config"
	"login-service/helper/at"
	"login-service/helper/atdb"
	"login-service/helper/watoken"
	"login-service/model"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDataBio(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	dockaryawan, err := atdb.GetOneDoc[model.Karyawan](config.Mongoconn, "kayawan", primitive.M{"nama": payload.Alias})
	if err != nil {
		dockaryawan.PhoneNumber = payload.Id
		dockaryawan.Nama = payload.Alias
		at.WriteJSON(respw, http.StatusOK, dockaryawan)
		return
	}
	dockaryawan.Nama = payload.Alias
	at.WriteJSON(respw, http.StatusOK, dockaryawan)
}