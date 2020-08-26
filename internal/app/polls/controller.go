package polls

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/julienschmidt/httprouter"
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
	pollID, errs := createPoll(poll)
	if errs != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(errs)
		return err
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"_id":%s}`, pollID)))
	return nil
}

// Show : GET "/polls/:id"
func (p *Polls) Show(w http.ResponseWriter, r *http.Request) error {
	ps := httprouter.ParamsFromContext(r.Context())
	_, err := primitive.ObjectIDFromHex(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		return err
	}
	poll := new(Poll)
	// TODO :
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// err = poll.Collection().FindOne(ctx, bson.D{{"_id", id}}).Decode(&poll)
	// if err != nil {
	// 	w.WriteHeader(404)
	// 	return
	// }
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
	pollID := ps.ByName("id")
	poll := new(PollParams)
	err := json.NewDecoder(r.Body).Decode(&poll)
	if err != nil {
		w.WriteHeader(400)
		return err
	}
	err = updatePoll(pollID, poll)
	if err != nil {
		w.WriteHeader(400)
		return err
	}
	return nil
}
