package main

import (
	"log"
	"net/http"

	"./controllers"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/newuser", controllers.CreateNewUser).Methods("POST")
	r.HandleFunc("/signin", controllers.Signin).Methods("POST")
	r.HandleFunc("/refresh", controllers.Refresh).Methods("POST")
	r.HandleFunc("/deleteAll", controllers.DeleteAllTokenUser).Methods("POST")
	r.HandleFunc("/delete", controllers.DeleteTokenUser).Methods("POST")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3000", nil))

}
