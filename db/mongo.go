package db

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tx-record/config"
)

func InitDB() {
	dbSource := fmt.Sprintf("mongodb://%s:%s", config.Config.Mongo.Host, config.Config.Mongo.Port)
	err := mgm.SetDefaultConfig(nil, "blockchain", options.Client().ApplyURI(dbSource))
	if err != nil {
		panic(" db connect failed.")
	}
}
