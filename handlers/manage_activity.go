package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"finder_api/models"
	"finder_api/utils"
)

func (h *Handler) CreateActivity(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	userDetails, err := h.DB.ValidateToken(authToken)

	if err != nil {
		log.Fatal(err)
	}

	user_id := fmt.Sprint(userDetails["user_id"])

	// Read body of request, but limit input to save server resources
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Fatal(err)
	}

	// Decode the json body into a User struct
	var activity models.Activity
	if err := json.Unmarshal(body, &activity); err != nil {
		// If error, return an HTTP StatusNotAcceptable
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, err.Error())
		fmt.Fprintf(w, string(body))

		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Create the message on the database
	if err = h.DB.CreateActivity(user_id, activity.Name, activity.Description, activity.Category, activity.Longitude, activity.Latitude); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (h *Handler) UpdateActivity(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	_, err := h.DB.ValidateToken(authToken)

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

	// Decode the json body into a User struct
	var activity models.Activity
	if err := json.Unmarshal(body, &activity); err != nil {
		// If error, return an HTTP StatusNotAcceptable
		w.WriteHeader(http.StatusNotAcceptable)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Create the message on the database
	if err = h.DB.UpdateActivity(activity.Activity_id,
		utils.NewNullString(activity.Name),
		utils.NewNullString(activity.Description),
		utils.NewNullString(activity.Category),
		utils.NewNullString(utils.FloatToString(activity.Longitude)),
		utils.NewNullString(utils.FloatToString(activity.Latitude))); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (h *Handler) GetAllUsersByActivityID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	_, err := h.DB.ValidateToken(authToken)

	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()

	// Read req parameter (activity_id)
	activity_id := r.FormValue("activity_id")

	if activity_id != "" {
		users, err := h.DB.RetrieveUsers(activity_id)
		if err != nil {
			fmt.Println(w, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Fatal(err)
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (h *Handler) DeleteActivity(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

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

		err = h.DB.DeleteActivity(activity_id, user_id)

		if err != nil {
			fmt.Println(w, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (h *Handler) GetActivityByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	_, err := h.DB.ValidateToken(authToken)

	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()

	// Read req parameter (activity_id)
	activity_id := r.FormValue("activity_id")

	activity, err := h.DB.GetActivityByID(activity_id)
	if err != nil {
		fmt.Println(w, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	marshaledActivity, err := json.Marshal(activity)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(marshaledActivity))
}

func (h *Handler) GetActivities(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	_, err := h.DB.ValidateToken(authToken)

	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()

	person_coord := new(utils.Coord)
	// Read req parameter (activity_id)
	category_filter := r.FormValue("category_filter")

	lat := 0.
	lon := 0.
	radius_filter, err := strconv.ParseFloat(r.FormValue("radius_filter"), 8)
	fmt.Println(radius_filter)
	if err != nil {
		radius_filter = 0.
	} else {
		lat, err = strconv.ParseFloat(r.FormValue("latitude"), 8)
		fmt.Println(lat)
		lon, err = strconv.ParseFloat(r.FormValue("longitude"), 8)
		fmt.Println(lon)
		person_coord = &utils.Coord{
			Lat: lat,
			Lon: lon,
		}
	}

	activities, err := h.DB.GetActivities(category_filter, radius_filter, *person_coord)
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
