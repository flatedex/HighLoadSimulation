package main

import(
	//"fmt"
	"net/http"
)

type Handler struct {}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request){
	_, err := writer.Write([]byte(http.StatusText(200)))
	if(err != nil) { panic(err) }
} 

func main(){
	http.ListenAndServe("127.0.0.1:8069", &Handler{})
}