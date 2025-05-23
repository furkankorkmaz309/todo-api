package handlers

import (
	"fmt"
	"net/http"
)

func WelcomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method!", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Welcome to todo-api!")
}
