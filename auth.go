package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "sign_up.html", nil)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicaton/json")
	var userResponse UserResponse
	r.ParseForm()
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")
	role := strings.ToLower(r.FormValue("role"))

	userResp, err := findUser(dbFindUrl, email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(userResp.Body) != 0 {
		fmt.Println("Email Already Taken, Please Choose Another one.")
		return
	}
	hashedPassword, hashedErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashedErr != nil {
		fmt.Println(hashedErr)
		return
	}

	var user = User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		Doctype:  "user",
	}

	jUser, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(jUser))
	request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jUser))
	if err != nil {
		fmt.Println("Byte Error", err)
		return
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		fmt.Println("Req Err", error)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		log.Fatalln("UnMarshal Err: ", error)
		return
	}
	userResp, err = findUser(dbFindUrl, email)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(userResp.Body) == 0 {
		fmt.Println("No User Data Found")
		return
	}
	user = userResp.Body[0]
	fmt.Println(user)
	if user.Role == "admin" {
		adminProfile := AdminProfile{
			// ID:      uuid.NewString(),
			// UserId:  user.ID,
			// Doctype: "adminProfile",
		}
		jData, err := json.Marshal(adminProfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jData))
		if err != nil {
			fmt.Println(err)
			return
		}
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(body)
		http.Redirect(w, r, "/sign_in", http.StatusPermanentRedirect)
		return
	} else if user.Role == "user" {
		userProfile := UserProfile{
			// ID:      uuid.NewString(),
			// UserId:  user.ID,
			// Doctype: "lecturerProfile",
		}
		jData, err := json.Marshal(userProfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jData))
		if err != nil {
			fmt.Println(err)
			return
		}
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(body)
		http.Redirect(w, r, "/sign_in", http.StatusPermanentRedirect)
		return
	} else if user.Role == "super_admin" {
		studentProfile := SuperAdminProfile{
			// ID:      uuid.NewString(),
			// UserId:  user.ID,
			// Doctype: "studentProfile",
		}
		jData, err := json.Marshal(studentProfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jData))
		if err != nil {
			fmt.Println(err)
			return
		}
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(body)
		http.Redirect(w, r, "/sign_in", http.StatusPermanentRedirect)
		return
	}

	http.Redirect(w, r, "/sign_in", http.StatusPermanentRedirect)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "sign_in.html", nil)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	r.ParseForm()
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	userResponse, err := findUser(dbFindUrl, email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(userResponse.Body) == 0 {
		fmt.Println("Invalid Username or Password")
		http.Redirect(w, r, "/sign_in", http.StatusPermanentRedirect)
		return
	}

	user = userResponse.Body[0]

	hash_err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if hash_err != nil {
		tmpl.ExecuteTemplate(w, "sign_in.html", "Please Verify Login Credentials")
		return
	}
	currentUser = user

	session, session_err := store.Get(r, "kaccima_session")
	if session_err != nil {
		http.Error(w, session_err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	session.Values["email"] = user.Email
	session.Values["_id"] = user.ID
	session_save_err := session.Save(r, w)
	if session_save_err != nil {
		fmt.Println(session_save_err)
		http.Error(w, session_save_err.Error(), http.StatusInternalServerError)
		return
	}
	//
	if user.Role == "admin" {
		http.Redirect(w, r, "/admin", http.StatusPermanentRedirect)
		return
	} else if user.Role == "student" {
		http.Redirect(w, r, "/student", http.StatusPermanentRedirect)
		return
	} else if user.Role == "lecturer" {
		http.Redirect(w, r, "/lecturer", http.StatusPermanentRedirect)
		return
	}
}

func signOut(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "kaccima_session")
	delete(session.Values, "email")
	delete(session.Values, "id")
	session.Save(r, w)
	currentUser = User{}
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func isAuthenticated(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "kaccima_session")
		email, ok := session.Values["email"]
		if !ok {
			http.Redirect(w, r, "/sign_in", http.StatusPermanentRedirect)
			return
		}

		userResp, err := findUser(dbFindUrl, fmt.Sprint(email))
		if err != nil {
			fmt.Println(err)
			return
		}
		currentUser = userResp.Body[0]
		endpoint(w, r)
	})
}
