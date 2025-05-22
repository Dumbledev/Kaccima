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

	http.HandleFunc("/", index)
	http.HandleFunc("/sign_up", signUp)
	http.HandleFunc("/sign_up_handler", signUpHandler)

	http.HandleFunc("/sign_in", signIn)
	http.HandleFunc("/sign_in_handler", signInHandler)

	http.HandleFunc("/sign_out", signOut)

	http.Handle("/profile", isAuthenticated(profile))
	http.Handle("/profile_update", isAuthenticated(profileUpdate))
	http.Handle("/profile_update_handler", isAuthenticated(profileUpdateHandler))

	fmt.Println("Server Running on Port :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func profile(w http.ResponseWriter, r *http.Request) {
	var profile UserProfile
	profileResp, err := findUserProfile(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	profile = profileResp.Body[0]
	tmpl.ExecuteTemplate(w, "profile.html", profile)
}

func profileUpdate(w http.ResponseWriter, r *http.Request) {
	var profile UserProfile
	profileResp, err := findUserProfile(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	profile = profileResp.Body[0]
	tmpl.ExecuteTemplate(w, "profile_update.html", profile)
}

func profileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// profileId := r.FormValue("profileId")
	http.Redirect(w, r, "/profile", http.StatusPermanentRedirect)
}
