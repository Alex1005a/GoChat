package Repositories

import (
	"awesomeProject1/src/app/domain"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("chat").Collection("users")
}

type (
	UserRepo struct {
	}
)

func NewUserRepo() domain.UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) CreateUser(name string, password string) (string, error) {
	var user domain.User
	user.Id = bson.NewObjectId().Hex()
	user.Name = name
	user.PasswordHash = password

	collection.InsertOne(ctx, user)

	return user.Id, nil
}

func (r *UserRepo) UserByID(id string) (domain.User, error) {
	var user domain.User

	_ = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return user, nil
}

func (r *UserRepo) UserByUsername(name string) (domain.User, error) {
	var user domain.User
	_ = collection.FindOne(ctx, bson.M{"name": name}).Decode(&user)

	return user, nil
}

func (r *UserRepo) DeleteUserByID(id string) error {
	collection.DeleteOne(ctx, bson.M{"_id": id})
	return  nil
}
