package controller

import (
	"fmt"
	"net/http"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "success")
}
