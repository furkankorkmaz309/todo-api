package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func takeIDFromURL(r *http.Request, index int) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) <= index {
		return 0, fmt.Errorf("index not found")
	}
	id, err := strconv.Atoi(parts[index])
	if err != nil {
		return 0, fmt.Errorf("an error occurred while converting string to integer: %v", err)
	}
	return id, nil
}
