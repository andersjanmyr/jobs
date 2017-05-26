package main

import "net/http"

func main() {
	http.HandleFunc("/", pipeline)
	http.ListenAndServe(":8080", nil)
}

func pipeline(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pipeline"))
}
