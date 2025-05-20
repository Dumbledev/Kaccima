package main

import (
	"fmt"
	"log"
	"net/http"

	"html/template"

	"github.com/gorilla/sessions"
)

var dbUrl = "http://admin:admin@localhost:5984/kaccima"
var dbFindUrl = dbUrl + "/_find"
var store = sessions.NewCookieStore([]byte("Secret"))
var currentUser User
var tmpl *template.Template

func main() {
	tmpl, _ = template.ParseGlob("templates/*.html")

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	http.HandleFunc("/sign_up", signUp)
	http.HandleFunc("/sign_up_handler", signUpHandler)

	http.HandleFunc("/sign_in", signIn)
	http.HandleFunc("/sign_in_handler", signInHandler)

	fmt.Println("Server Running on Port :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
