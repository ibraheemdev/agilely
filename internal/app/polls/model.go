package polls

import (
	"github.com/google/uuid"
	"github.com/ibraheemdev/agilely/pkg/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Poll Document
type Poll struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Password uuid.UUID          `json:"password" bson:"password"`
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
	// TODO :
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// pwd, err := uuid.NewRandom()
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return "", nil
	// }
	// create := bson.M{"title": poll.Title, "password": pwd}
	// res, err := config.DatabaseClient.Collection("polls").InsertOne(ctx, create)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return "", nil
	// }
	// id := res.InsertedID.(primitive.ObjectID).Hex()
	// return id, nil
	return "", nil
}

func updatePoll(id string, poll *PollParams) error {
	// TODO :
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// filter := bson.D{{"_id", id}}
	// update := bson.D{{"$set", bson.D{{"title", poll.Title}}}}
	// // _, err := config.DatabaseClient.Collection("polls").UpdateOne(ctx, filter, update)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return err
	// }
	return nil
}
