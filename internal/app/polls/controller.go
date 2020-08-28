package polls

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// Collection ...
	Collection = "polls"
)

// Polls Controller
type Polls struct {
	*engine.Engine
}

// NewController : Returns a new polls controller
func NewController(e *engine.Engine) *Polls {
	return &Polls{Engine: e}
}

// Create : POST "/polls"
func (p *Polls) Create(w http.ResponseWriter, r *http.Request) error {
	poll := new(PollParams)
	err := json.NewDecoder(r.Body).Decode(&poll)
	if err != nil {
		w.WriteHeader(400)
		return err
	}
	errs := validate(poll)
	if errs != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(errs)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pwd, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return err
	}

	create := bson.M{"title": poll.Title, "password": pwd}
	res, err := p.Core.Database.Collection("polls").InsertOne(ctx, create)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"_id": %s}`, id)))
	return nil
}

// Show : GET "/polls/:id"
func (p *Polls) Show(w http.ResponseWriter, r *http.Request) error {
	ps := httprouter.ParamsFromContext(r.Context())
	id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		return err
	}
	poll := new(Poll)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = p.Core.Database.Collection(Collection).FindOne(ctx, bson.D{{"_id", id}}).Decode(&poll)
	if err != nil {
		w.WriteHeader(404)
		return nil
	}
	json, err := json.Marshal(poll)
	if err != nil {
		w.WriteHeader(400)
		return err
	}
	if _, err = w.Write(json); err != nil {
		log.Fatal(err)
	}
	return nil
}

// Update : PUT "/polls/:id"
func (p *Polls) Update(w http.ResponseWriter, r *http.Request) error {
	ps := httprouter.ParamsFromContext(r.Context())
	id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	poll := new(PollParams)
	err = json.NewDecoder(r.Body).Decode(&poll)
	if err != nil {
		w.WriteHeader(400)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{"_id", id}}

	update := bson.D{{"$set", bson.D{{"title", poll.Title}}}}
	_, err = p.Core.Database.Collection(Collection).UpdateOne(ctx, filter, update)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Println(err.Error())
		return err
	}

	return nil
}
