package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"

	"finder_api/models"
)

// POST "/signup"
func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Read body of request, but limit input to save server resources
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Fatal(err)
	}

	// Decode the json body into a User struct
	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		// If error, return an HTTP StatusNotAcceptable
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, err.Error())
		fmt.Fprintf(w, string(body))

		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	// Create the message on the database
	if err = h.DB.UserSignup(user.First_name, user.Last_name, user.Email, user.Phone_number, hashedPassword); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
