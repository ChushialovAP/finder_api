package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"finder_api/models"
)

// POST "/auth"
func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Read body of request, but limit input to save server resources
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Fatal(err)
	}

	// Unmarshal the json body into a User struct
	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		// If error, return an HTTP StatusNotAcceptable
		w.WriteHeader(http.StatusNotAcceptable)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	tokenDetails, err := h.DB.GenerateToken(user.Phone_number, user.Password)

	if err != nil {
		log.Fatal(err)
	}

	marshaledTokenDetails, err := json.Marshal(tokenDetails)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(marshaledTokenDetails))
}
