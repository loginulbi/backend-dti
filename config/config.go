package config

import (
	"log"
	"os"

	"login-service/helper/at"
	"login-service/helper/atdb"
	"login-service/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var PrivateKey string = os.Getenv("PRKEY")

var IPPort, Net = at.GetAddress()

var PhoneNumber string = os.Getenv("PHONENUMBER")

func SetEnv() {
	if ErrorMongoconn != nil {
		log.Println(ErrorMongoconn.Error())
	}
	profile, err := atdb.GetOneDoc[model.Profile](Mongoconn, "profile", primitive.M{"phonenumber": PhoneNumber})
	if err != nil {
		log.Println(err)
	}
	if Mongoconn == nil {
		log.Println("Failed to connect to MongoDB.")
		return
	}
	PublicKeyWhatsAuth = profile.PublicKey
	WAAPIToken = profile.Token
}
