package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"finder_api/models"
)

// GET "/messages"
func (h *Handler) MessagesGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	_, err := h.DB.ValidateToken(authToken)

	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()

	// Read req parameter (activity_id)
	activity_id := r.FormValue("activity_id")

	messages, err := h.DB.MessagesRetrieve(activity_id)
	if err != nil {
		fmt.Println(w, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Fatal(err)
	}
}

// POST "/messages"
func (h *Handler) MessagesPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	userDetails, err := h.DB.ValidateToken(authToken)

	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()

	user_id := fmt.Sprint(userDetails["user_id"])
	//fmt.Println(w, user_id)

	// Read body of request, but limit input to save server resources
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Fatal(err)
	}

	// Unmarshal the json body into a Message struct
	var message models.Message
	if err := json.Unmarshal(body, &message); err != nil {
		// If error, return an HTTP StatusNotAcceptable
		w.WriteHeader(http.StatusNotAcceptable)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Create the message on the database
	if err = h.DB.MessageCreate(user_id, message.Activity_id, message.Text); err != nil {
		fmt.Println(w, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
