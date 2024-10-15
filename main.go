package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/infocenter/", infocenterHandler)
	http.ListenAndServe(":8080", nil)
}

func infocenterHandler(w http.ResponseWriter, r *http.Request) {

}
