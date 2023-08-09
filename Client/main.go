package main

import (
	"fmt"
	"net/http"
)
const PORT int = 8069
type msg string

func (m msg) ServeHTTP(resp http.ResponseWriter, req *http.Request){
	fmt.Fprint(resp, m)
}

func main() { 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "Static/home.html")
	})

	fmt.Println("Server is listening at localhost:8069")
	http.ListenAndServe(":8069", http.FileServer(http.Dir("static")))
}