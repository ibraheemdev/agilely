package polls

import (
	"github.com/ibraheemdev/agilely/pkg/validator"

	"github.com/google/uuid"
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
