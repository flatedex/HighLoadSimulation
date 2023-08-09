package main

import (
	_ "fmt"
	"net/http"
)

type Message struct {
	Name string
	Message string
}

func GetDataFromClient(writer http.ResponseWriter, request *http.Request){
	defer writer.Write([]byte(http.StatusText(200)))

	err := request.ParseForm()
	if(err != nil) { panic(err) }

	var message Message

	message.Name = request.FormValue("Name")
	message.Message = request.FormValue("Message")
}

func main(){	
	http.HandleFunc("/home", GetDataFromClient)
	http.ListenAndServe("127.0.0.1:8069", nil)
}