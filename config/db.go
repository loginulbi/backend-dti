package config

import (
	"login-service/helper"
	"login-service/helper/atdb"
	"os"
)

var MongoString string = os.Getenv("MONGOSTRING")

var mongoinfo = atdb.DBInfo{
	DBString: helper.SRVLookup(MongoString),
	DBName:   "hris",
}

var Mongoconn, _ = atdb.MongoConnect(mongoinfo)
