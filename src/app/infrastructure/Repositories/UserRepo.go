package Repositories

import (
	"awesomeProject1/src/app/domain"
	"awesomeProject1/src/app/infrastructure"
	"context"
	"github.com/dgrijalva/jwt-go"
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
	hash := hasher.NewHasher()
	result, _ := hash.PasswordToHash(password)
	user.PasswordHash = result

	collection.InsertOne(ctx, user)

	return user.Id, nil
}

func (r *UserRepo) UserByID(id string) (domain.User, error) {
	var user domain.User

	_ = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	user.PasswordHash = ""

	return user, nil
}

func (r *UserRepo) UserByUsername(name string) (domain.User, error) {
	var user domain.User
	_ = collection.FindOne(ctx, bson.M{"name": name}).Decode(&user)
	user.PasswordHash = ""

	return user, nil
}

func (r *UserRepo) DeleteUserByID(id string) error {
	collection.DeleteOne(ctx, bson.M{"_id": id})
	return nil
}

func (r *UserRepo) Login(name, password string) domain.User {

	var user domain.User
	_ = collection.FindOne(ctx, bson.M{"name": name}).Decode(&user)

	hash := hasher.NewHasher()
	result := hash.CheckPassword(password, user.PasswordHash)
	if result == false {
		return domain.User{}
	}

	user.PasswordHash = ""

	tk := &domain.Token{UserId: user.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("token_password"))
	user.Token = tokenString

	//var resp map[string]interface{}
	//resp["account"] = user
	return user
}
