package handlers

import (
	"fmt"

	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) TestResourceHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	userDetails, err := h.DB.ValidateToken(authToken)

	if err != nil {

		fmt.Fprintf(w, err.Error())

	} else {

		name := fmt.Sprint(userDetails["first_name"])
		id := fmt.Sprint(userDetails["user_id"])

		fmt.Fprintf(w, "Welcome, "+name+"\r\n"+id)
	}

}
