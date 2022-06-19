package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) JoinActivity(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	userDetails, err := h.DB.ValidateToken(authToken)

	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()

	// Read req parameter (activity_id)
	activity_id := r.FormValue("activity_id")

	if activity_id != "" {
		user_id := fmt.Sprint(userDetails["user_id"])
		err := h.DB.JoinActivity(user_id, activity_id)

		if err != nil {
			fmt.Fprintf(w, "user has already joined the activity")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "user joined activity with activity_id = "+activity_id)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}
