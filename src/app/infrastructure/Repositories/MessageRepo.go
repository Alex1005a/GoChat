package Repositories

import (
	"awesomeProject1/src/app/domain"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

var messageCollection *mongo.Collection
var messagectx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(messagectx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(messagectx, nil)
	if err != nil {
		log.Fatal(err)
	}

	messageCollection = client.Database("chat").Collection("messages")
}

type (
	MessageRepo struct {
	}
)

func NewMessageRepo() domain.MessageRepo {
	return &MessageRepo{}
}

func (r *MessageRepo) CreateMessage(text string) (domain.Message, error) {
	var mes domain.Message
	mes.Id = bson.NewObjectId().Hex()
	mes.Text = text
	mes.CreatedTime = time.Now()

	messageCollection.InsertOne(messagectx, mes)

	return mes, nil
}

func (r *MessageRepo) GetFiveLastMessages() []domain.Message{
	var messages []domain.Message
	options := options.Find()

	//options.SetSort(bson.D{{"createdTime", -1}})

	options.SetLimit(5)
	cur, _ := messageCollection.Find(messagectx, bson.D{}, options)
	for cur.Next(context.Background()) {

		var l domain.Message
		_ = cur.Decode(&l)

		messages = append(messages, l)
	}
	return messages
}

