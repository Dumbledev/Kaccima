package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "io"
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
	r.ParseMultipartForm(10 * 1024 * 1024)
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	userResponse, err := findUserByEmail(dbFindUrl, email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(userResponse.Body) != 0 {
		tmpl.ExecuteTemplate(w, "sign_up.html", "User Already Registered, Please Register With Other Credentials.")
		return
	}
	hashedPassword, hashedErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashedErr != nil {
		tmpl.ExecuteTemplate(w, "sign_up.html", "Error")
		return
	}

	var org = User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
		Doctype:  "user",
	}

	jsonData, err := json.Marshal(org)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tmpl.ExecuteTemplate(w, "sign_up.html", "Server Error(1)")
		return
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		tmpl.ExecuteTemplate(w, "sign_up.html", "Server Error(2)")
	}
	defer res.Body.Close()
	http.Redirect(w, r, "/sign_in", http.StatusPermanentRedirect)
}

func adminSignUp(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin_sign_up.html", nil)
}

func adminSignUpHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	userResponse, err := findUserByEmail(dbFindUrl, email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(userResponse.Body) != 0 {
		tmpl.ExecuteTemplate(w, "admin_sign_in.html", "Email Already Registered, Please Choose Another one.")
		return
	}
	hashedPassword, hashedErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashedErr != nil {
		tmpl.ExecuteTemplate(w, "admin_sign_up.html", "Error")
		return
	}

	var org = User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hashedPassword),
		Role:     "admin",
		Doctype:  "user",
	}

	jsonData, err := json.Marshal(org)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(jUser))
	request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tmpl.ExecuteTemplate(w, "admin_sign_up.html", "Server Error(1)")
		return
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		tmpl.ExecuteTemplate(w, "admin_sign_up.html", "Server Error(2)")
		return
	}
	defer res.Body.Close()

	http.Redirect(w, r, "/sign_in", http.StatusPermanentRedirect)
}

func superAdminSignUp(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "super_admin_signup.html", nil)
}

func superAdminSignUpHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	userResponse, err := findUserByEmail(dbFindUrl, email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(userResponse.Body) != 0 {
		tmpl.ExecuteTemplate(w, "super_admin_signup.html", "Email Already Registered, Please Choose Another one.")
		return
	}
	hashedPassword, hashedErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashedErr != nil {
		tmpl.ExecuteTemplate(w, "super_admin_signup.html", "Error")
		return
	}

	var org = User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hashedPassword),
		Role:     "superAdmin",
		Doctype:  "user",
	}

	jsonData, err := json.Marshal(org)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(jUser))
	request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tmpl.ExecuteTemplate(w, "super_admin_signup.html", "Server Error(1)")
		return
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		tmpl.ExecuteTemplate(w, "super_admin_signup.html", "Server Error(2)")
		return
	}
	defer res.Body.Close()

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

	userResponse, err := findUserByEmail(dbFindUrl, email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(userResponse.Body) == 0 {
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
		http.Error(w, session_save_err.Error(), http.StatusInternalServerError)
		return
	}
	//
	if user.Role == "admin" {
		http.Redirect(w, r, "/admin", http.StatusPermanentRedirect)
		return
	} else if user.Role == "user" {
		http.Redirect(w, r, "/dashboard", http.StatusPermanentRedirect)
		return
	} else if user.Role == "superAdmin" {
		http.Redirect(w, r, "/approval_admin", http.StatusSeeOther)
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

func forgotPassword(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "forgot-password.html", nil)
}

func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(r.FormValue("email"))
	userResp, err := findUserByEmail(dbFindUrl, email)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(userResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "forgot-password.html", "Invalid or Wrong Email")
		return
	}
	user := userResp.Body[0]
	http.Redirect(w, r, "/password_reset_qa/"+user.Email, http.StatusSeeOther)
}

func passwordResetQa(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(r.PathValue("email"))
	tmpl.ExecuteTemplate(w, "password_reset_qa.html", email)
}

func passwordResetQaHandler(w http.ResponseWriter, r *http.Request) {
	var securityQa PasswordResetQuestion
	r.ParseForm()
	question := r.FormValue("security_question")
	answer := strings.ToLower(r.FormValue("security_answer"))
	email := strings.ToLower(r.FormValue("email"))
	fmt.Println(email, question, answer)
	userResp, err := findUserByEmail(dbFindUrl, email)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(userResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "forgot-password.html", "Invalid or Wrong Email")
		return
	}
	user := userResp.Body[0]
	secQaResp, err := findSecurityQAByUser(dbFindUrl, user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(secQaResp.Body) == 0 {
		securityQa = PasswordResetQuestion{}
	} else {
		securityQa = secQaResp.Body[0]
	}
	if question == securityQa.Question && answer == securityQa.Answer {
		http.Redirect(w, r, "/reset_password/"+email, http.StatusSeeOther)
	} else {
		tmpl.ExecuteTemplate(w, "forgot-password.html", "Try Again")
	}
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	tmpl.ExecuteTemplate(w, "reset_password.html", email)
}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	newPassword := r.FormValue("newPassword")
	userResp, err := findUserByEmail(dbFindUrl, email)
	if err != nil {
		fmt.Println(err)
		return
	}
	user := userResp.Body[0]

	hashedPassword, hashedErr := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if hashedErr != nil {
		tmpl.ExecuteTemplate(w, "reset_password.html", "Error Adding New Password")
		return
	}
	user.Password = string(hashedPassword)
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+user.ID, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header.Set("content-type", "application/json")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))
	http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
}

func isAuthenticated(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "kaccima_session")
		email, ok := session.Values["email"]
		if !ok {
			http.Redirect(w, r, "/sign_in", http.StatusPermanentRedirect)
			return
		}

		userResp, err := findUserByEmail(dbFindUrl, fmt.Sprint(email))
		if err != nil {
			fmt.Println(err)
			return
		}
		currentUser = userResp.Body[0]
		endpoint(w, r)
	})
}
