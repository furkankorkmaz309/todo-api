package handlers

import (
	"net/http"
)

func WelcomePage(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, nil, "Welcome to todo-api!")
}
