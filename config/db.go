package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DB KARYAWAN
func DBPresensi(dbname string) *mongo.Database {
	connectionstr := os.Getenv("KARYAWANDATA")
	if connectionstr == "" {
		panic(fmt.Errorf("KARYAWANDATA ENV NOT FOUND"))
	}
	clay, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionstr))
	if err != nil {
		panic(fmt.Errorf("MongoConnect: %+v \n", err))
	}
	return clay.Database(dbname)
}




// var MongoString string = os.Getenv("MONGOSTRING")

// var mongoinfo = atdb.DBInfo{
// 	DBString: MongoString,
// 	DBName:   "hris",
// }

// var Mongoconn, _ = atdb.MongoConnect(mongoinfo)
