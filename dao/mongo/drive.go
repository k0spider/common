package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoConfig struct {
	AppName     string        `yaml:"appName"`
	Database    string        `yaml:"database"`
	Dns         string        `yaml:"dns"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
	MaxOpens    uint64        `yaml:"maxOpens"`
	MinOpens    uint64        `yaml:"minOpens"`
}

type MongoDrive struct {
	*mongo.Database
}

type MongoClient struct {
	c *MongoConfig
	*mongo.Client
}

func NewMongo(c *MongoConfig) *MongoDrive {
	ctx := context.Background()
	db := new(MongoDrive)
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(c.Dns).
		SetAppName(c.AppName).
		SetMaxConnIdleTime(time.Millisecond*c.IdleTimeout).
		SetMaxPoolSize(c.MaxOpens).
		SetMinPoolSize(c.MinOpens))
	if err != nil {
		panic(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	db.Database = client.Database(c.Database)
	return db
}

func NewMongoClient(c *MongoConfig) *MongoClient {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(c.Dns).
		SetAppName(c.AppName).
		SetMaxConnIdleTime(time.Millisecond*c.IdleTimeout).
		SetMaxPoolSize(c.MaxOpens).
		SetMinPoolSize(c.MinOpens))
	if err != nil {
		panic(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	mc := &MongoClient{
		c:      c,
		Client: client,
	}
	return mc
}

func (this *MongoClient) GetConn() (mongo.Session, *mongo.Database) {
	session, err := this.StartSession()
	if err != nil {
		panic(err)
	}

	return session, session.Client().Database(this.c.Database)
}

func (this *MongoClient) CloseConn(session mongo.Session, c ...context.Context) {
	var ctx context.Context
	if len(c) > 0 {
		ctx = c[0]
	} else {
		ctx = context.Background()
	}
	session.EndSession(ctx)
}
