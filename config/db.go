package config

import (
	"os"

	"login-service/helper/atdb"
)

var MongoString string = os.Getenv("MONGOSTRING")

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "hris",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
