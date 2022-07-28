package db

import (
	"context"
	"demo/db/model"
	"fmt"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoInstance *mongo.Client

const DataBase = "cloud"

func InitMongoDB() error {
	user := os.Getenv("MONGO_USERNAME")
	pwd := os.Getenv("MONGO_PASSWORD")
	mongoAddress := os.Getenv("MONGO_ADDRESS")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tmp, err := url.Parse(mongoAddress)
	if err != nil {
		fmt.Printf("mongoAddress parse error err %v", err)
		return err
	}
	authSource := tmp.Query().Get("authSource")

	credential := options.Credential{
		AuthSource: authSource,
		Username:   user,
		Password:   pwd,
	}

	mongoUrl := fmt.Sprintf("mongodb://%s", mongoAddress)

	clientOpts := options.Client().ApplyURI(mongoUrl).SetAuth(credential)

	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		return err
	}

	coll := client.Database(DataBase).Collection("count")
	doc := &model.MongoCount{
		Type:  "mongodb",
		Count: 2022,
	}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("documents inserted with result:%v\n", result)

	mongoInstance = client
	return err
}

func GetMongo() *mongo.Client {
	return mongoInstance
}
