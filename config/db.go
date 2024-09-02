package config

import (
	"login-service/helper/atdb"
	"os"
)

var MongoString string = os.Getenv("MONGOSTRING")

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "hris",
}

var Mongoconn, _ = atdb.MongoConnect(mongoinfo)
