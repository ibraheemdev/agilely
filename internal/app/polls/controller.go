package polls

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Create : POST "/polls"
func Create() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		poll := new(PollParams)
		err := json.NewDecoder(r.Body).Decode(&poll)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		pollID, errs := createPoll(poll)
		if errs != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(errs)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"_id":%s}`, pollID)))
	}
}

// Show : GET "/polls/:id"
func Show() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			w.WriteHeader(400)
			return
		}
		poll := new(Poll)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = poll.Collection().FindOne(ctx, bson.D{{"_id", id}}).Decode(&poll)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		json, err := json.Marshal(poll)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		if _, err = w.Write(json); err != nil {
			log.Fatal(err)
		}
	}
}

// Update : PUT "/polls/:id"
func Update() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		pollID := ps.ByName("id")
		poll := new(PollParams)
		err := json.NewDecoder(r.Body).Decode(&poll)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		err = updatePoll(pollID, poll)
		if err != nil {
			w.WriteHeader(400)
			return
		}
	}
}
