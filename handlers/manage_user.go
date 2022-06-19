package handlers

import (
	"encoding/json"
	"finder_api/models"
	"finder_api/utils"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) GetAllActivitiesByUserID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	userDetails, err := h.DB.ValidateToken(authToken)

	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()

	user_id := fmt.Sprint(userDetails["user_id"])

	activities, err := h.DB.RetriveActivities(user_id)
	if err != nil {
		fmt.Println(w, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(activities); err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	userDetails, err := h.DB.ValidateToken(authToken)

	if err != nil {
		log.Fatal(err)
	}

	// Read body of request, but limit input to save server resources
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Fatal(err)
	}

	user_id := fmt.Sprint(userDetails["user_id"])

	// Decode the json body into a User struct
	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		// If error, return an HTTP StatusNotAcceptable
		w.WriteHeader(http.StatusNotAcceptable)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Create the message on the database
	if err = h.DB.UpdateUser(user_id,
		utils.NewNullString(user.First_name),
		utils.NewNullString(user.Last_name),
		utils.NewNullString(user.Email),
		utils.NewNullString(user.Phone_number)); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
