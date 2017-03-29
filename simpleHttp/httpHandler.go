package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// cookie handling



// login handler
func homePageHandler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//w.Header().Set("content-type", "application/json") 
	fmt.Fprint(w, "Hello World")
}



//home page handler

// server main method

var router = mux.NewRouter()

func main() {

	router.HandleFunc("/query", homePageHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}