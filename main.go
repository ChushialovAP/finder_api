package main

import (
	"net/http"

	"log"

	"github.com/julienschmidt/httprouter"

	//"github.com/todsul/chat-api-golang/handlers"
	//"github.com/todsul/chat-api-golang/models"

	"finder_api/handlers"
	"finder_api/models"
)

func main() {

	db := models.New()
	router := httprouter.New()
	handler := &handlers.Handler{db}

	// We wrap the handler to execute boilerplate handler code
	router.Handle("GET", "/messages", handler.Process(handler.MessagesGet))
	router.Handle("POST", "/messages", handler.Process(handler.MessagesPost))
	router.Handle("POST", "/signup", handler.Process(handler.SignupHandler))
	router.Handle("POST", "/auth", handler.Process(handler.AuthHandler))
	router.Handle("PUT", "/user", handler.Process(handler.UpdateUser))
	//router.Handle("POST", "/auth/refresh", handler.Process(handler.RefreshToken))
	router.Handle("GET", "/test", handler.Process(handler.TestResourceHandler))
	router.Handle("GET", "/activity", handler.Process(handler.GetActivityByID))
	router.Handle("POST", "/activity", handler.Process(handler.CreateActivity))
	router.Handle("PUT", "/activity", handler.Process(handler.UpdateActivity))
	router.Handle("GET", "/join", handler.Process(handler.JoinActivity))
	router.Handle("GET", "/getusers", handler.Process(handler.GetAllUsersByActivityID))
	router.Handle("GET", "/getactivities", handler.Process(handler.GetActivities))
	router.Handle("GET", "/getuseractivities", handler.Process(handler.GetAllActivitiesByUserID))
	router.Handle("GET", "/deleteactivity", handler.Process(handler.DeleteActivity))

	log.Fatal(http.ListenAndServe(":8080", router))
}
