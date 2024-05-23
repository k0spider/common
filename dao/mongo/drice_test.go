package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type Student struct {
	Id              primitive.ObjectID `bson:"_id"`
	AppId           int64              `bson:"appid"`
	Event           string             `bson:"event"`
	Addresses       string             `bson:"addresses"`
	time            int64              `bson:"time"`
	Cid             string             `bson:"cid"`
	Sn              string             `bson:"sn"`
	Ip              string             `bson:"ip"`
	Country         string             `bson:"country"`
	Region          string             `bson:"region"`
	DeviceType      int                `bson:"device_type"`
	PhoneBrand      string             `bson:"phone_brand"`
	SysVer          string             `bson:"sys_ver"`
	NetWork         int                `bson:"network"`
	symbol          string             `bson:"symbol"`
	Model           string             `bson:"model"`
	Brand           string             `bson:"brand"`
	DownloadChannel string             `bson:"download_channel"`
	gas             string             `bson:"gas"`
	Reward          string             `bson:"reward"`
	Balance         string             `bson:"balance"`
	tableName       string             `bson:"-"`
}

func TestNewMongo(t *testing.T) {
	c := &MongoConfig{
		AppName:     "xxxx-project",
		Database:    "test",
		Dns:         "mongodb://user:pwd@127.0.0.1:27017",
		IdleTimeout: 30,
		MaxOpens:    30,
		MinOpens:    5,
	}
	testDB := NewMongo(c)
	find, err := testDB.Collection("test").Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	res := []*Student{}
	err = find.All(context.Background(), &res)
	if err != nil {
		panic(err)
	}
	fmt.Println(res[0])
}
