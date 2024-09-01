package whatsauth

import (
	"strings"

	"login-service/helper/atapi"
	"login-service/helper/atdb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func WebHook(WAKeyword, WAPhoneNumber, WAAPIQRLogin, WAAPIMessage string, msg IteungMessage, db *mongo.Database) (resp Response, err error) {
	if IsLoginRequest(msg, WAKeyword) { //untuk whatsauth request login
		resp, err = HandlerQRLogin(msg, WAKeyword, WAPhoneNumber, db, WAAPIQRLogin)
	} else { //untuk membalas pesan masuk
		resp, err = HandlerIncomingMessage(msg, WAPhoneNumber, db, WAAPIMessage)
	}
	return
}

func RefreshToken(dt *WebHookInfo, WAPhoneNumber, WAAPIGetToken string, db *mongo.Database) (res *mongo.UpdateResult, err error) {
	profile, err := GetAppProfile(WAPhoneNumber, db)
	if err != nil {
		return
	}
	var resp User
	if profile.Token != "" {
		_, resp, err = atapi.PostStructWithToken[User]("Token", profile.Token, dt, WAAPIGetToken)
		if err != nil {
			return
		}
		profile.Phonenumber = resp.PhoneNumber
		profile.Token = resp.Token
		res, err = atdb.ReplaceOneDoc(db, "profile", bson.M{"phonenumber": resp.PhoneNumber}, profile)
		if err != nil {
			return
		}
	}
	return
}

func IsLoginRequest(msg IteungMessage, keyword string) bool {
	return strings.Contains(msg.Message, keyword) && msg.From_link
}

func GetUUID(msg IteungMessage, keyword string) string {
	return strings.Replace(msg.Message, keyword, "", 1)
}

func HandlerQRLogin(msg IteungMessage, WAKeyword string, WAPhoneNumber string, db *mongo.Database, WAAPIQRLogin string) (resp Response, err error) {
	dt := &WhatsauthRequest{
		Uuid:        GetUUID(msg, WAKeyword),
		Phonenumber: msg.Phone_number,
		Delay:       msg.From_link_delay,
	}
	structtoken, err := GetAppProfile(WAPhoneNumber, db)
	if err != nil {
		return
	}
	_, resp, err = atapi.PostStructWithToken[Response]("Token", structtoken.Token, dt, WAAPIQRLogin)
	return
}

func HandlerIncomingMessage(msg IteungMessage, WAPhoneNumber string, db *mongo.Database, WAAPIMessage string) (resp Response, err error) {
	dt := &TextMessage{
		To:       msg.Chat_number,
		IsGroup:  false,
		Messages: GetRandomReplyFromMongo(msg, db),
	}
	if msg.Chat_server == "g.us" { //jika pesan datang dari group maka balas ke group
		dt.IsGroup = true
	}
	if (msg.Phone_number != "628112000279") && (msg.Phone_number != "6283131895000") { //ignore pesan datang dari iteung
		var profile Profile
		profile, err = GetAppProfile(WAPhoneNumber, db)
		if err != nil {
			return
		}
		_, resp, err = atapi.PostStructWithToken[Response]("Token", profile.Token, dt, WAAPIMessage)
		if err != nil {
			return
		}
	}
	return
}

func GetRandomReplyFromMongo(msg IteungMessage, db *mongo.Database) string {
	rply, err := atdb.GetRandomDoc[Reply](db, "reply", 1)
	if err != nil {
		return "Koneksi Database Gagal: " + err.Error()
	}
	replymsg := strings.ReplaceAll(rply[0].Message, "#BOTNAME#", msg.Alias_name)
	replymsg = strings.ReplaceAll(replymsg, "\\n", "\n")
	return replymsg
}

func GetAppProfile(phonenumber string, db *mongo.Database) (apitoken Profile, err error) {
	filter := bson.M{"phonenumber": phonenumber}
	apitoken, err = atdb.GetOneDoc[Profile](db, "profile", filter)

	return
}
