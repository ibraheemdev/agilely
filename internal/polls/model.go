package polls

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ibraheemdev/poller/pkg/base"
	"github.com/ibraheemdev/poller/pkg/database"
	"github.com/ibraheemdev/poller/pkg/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Poll Document
type Poll struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	Password   uuid.UUID          `json:"password" bson:"password"`
	base.Model `json:"-"`
}

// Collection : The poll collection
func (p Poll) Collection() base.ModelCollection {
	return base.ModelCollection{database.Client.Collection("polls")}
}

// PollParams : Valid poll params
type PollParams struct {
	Title string `json:"title" bson:"title"`
}

func validate(poll *PollParams) validator.ValidationErrors {
	v := &validator.Validator{}
	v.ValidatePresenceOf("Title", poll.Title)
	if errs := v.Errors; errs != nil {
		return validator.Stringify(errs)
	}
	return nil
}

func createPoll(poll *PollParams) (string, validator.ValidationErrors) {
	errs := validate(poll)
	if errs != nil {
		return "", errs
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pwd, err := uuid.NewRandom()
	if err != nil {
		log.Println(err.Error())
		// TODO : Move to base
		return "", nil
	}
	create := bson.M{"title": poll.Title, "password": pwd}
	res, err := database.Client.Collection("polls").InsertOne(ctx, create)
	if err != nil {
		log.Println(err.Error())
		// TODO : Move to base
		return "", nil
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func updatePoll(id string, poll *PollParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"title", poll.Title}}}}
	_, err := database.Client.Collection("polls").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
