package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"finder_api/models"
)

type Handler struct {
	DB models.Datastore
}

// Generic method to handle common tasks such as logging
func (h *Handler) Process(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Start the stopwatch for logging the request
		start := time.Now()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Call the handler and pass the original parameters
		handle(w, r, p)

		// Print the log to the console
		log.Printf(
			"%-6s%-20s%-20s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	}
}
