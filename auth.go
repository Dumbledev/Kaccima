package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

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
	w.Header().Set("Content-Type", "applicaton/json")
	r.ParseMultipartForm(10 * 1024 * 1024)
	name := r.FormValue("companyName")
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")
	address := r.FormValue("officeAddress")
	employeesNo := r.FormValue("employeesNo")
	nonNigerianEmployees := r.FormValue("nonNigerianEmployees")
	directorsNo := r.FormValue("directorsNo")
	nonNigerianDirectors := r.FormValue("nonNigerianDirectors")
	businessNature := r.FormValue("businessNature")
	bankers := r.FormValue("bankers")
	contactPerson := r.FormValue("contactPerson")
	rep := r.FormValue("representative")

	coverLetter, _, err := r.FormFile("coverLetter")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer coverLetter.Close()
	temp1, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp1.Close()
	fileBytes, err3 := io.ReadAll(coverLetter)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp1.Write(fileBytes)

	memorandum, _, err := r.FormFile("memorandum")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer memorandum.Close()
	temp2, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp2.Close()
	fileBytes1, err3 := io.ReadAll(memorandum)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp2.Write(fileBytes1)

	businessCertificate, _, err := r.FormFile("registrationCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer businessCertificate.Close()
	temp3, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp3.Close()
	fileBytes2, err3 := io.ReadAll(businessCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp2.Write(fileBytes2)

	incorporationCertificate, _, err := r.FormFile("incorporationCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer incorporationCertificate.Close()
	temp4, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp4.Close()
	fileBytes3, err3 := io.ReadAll(incorporationCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp2.Write(fileBytes3)

	passportPhoto, _, err := r.FormFile("passportPhotos")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer passportPhoto.Close()
	temp5, err2 := os.CreateTemp("static", "file-*.png")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp5.Close()
	fileBytes4, err3 := io.ReadAll(passportPhoto)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp5.Write(fileBytes4)

	businessPremiseCertificate, _, err := r.FormFile("premisesCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer businessPremiseCertificate.Close()
	temp6, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp6.Close()
	fileBytes5, err3 := io.ReadAll(businessPremiseCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp6.Write(fileBytes5)

	formC07, _, err := r.FormFile("formC07")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer formC07.Close()
	temp7, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp7.Close()
	fileBytes6, err3 := io.ReadAll(formC07)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp7.Write(fileBytes6)

	idDocument, _, err := r.FormFile("nationalId")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer idDocument.Close()
	temp8, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp8.Close()
	fileBytes7, err3 := io.ReadAll(idDocument)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp5.Write(fileBytes7)

	orgResp, err := findOrganization(dbFindUrl, email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(orgResp.Body) != 0 {
		tmpl.ExecuteTemplate(w, "sign_in.html", "Organization Already Registered, Please Choose Another one.")
		return
	}
	hashedPassword, hashedErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashedErr != nil {
		tmpl.ExecuteTemplate(w, "sign_up.html", "Error")
		return
	}

	var org = Organization{
		ID:                         uuid.NewString(),
		Name:                       name,
		Address:                    address,
		Email:                      email,
		NumberOfEmployees:          employeesNo,
		NonNigerianEmployees:       nonNigerianEmployees,
		NumberOfDirectors:          directorsNo,
		NonNigerianDirectors:       nonNigerianDirectors,
		NatureOfBusiness:           businessNature,
		Bankers:                    bankers,
		ContactPerson:              contactPerson,
		Representative:             rep,
		Password:                   string(hashedPassword),
		CoverLetter:                temp1.Name(),
		Memorandum:                 temp2.Name(),
		BusinessCertificate:        temp3.Name(),
		IncorporationCertificate:   temp4.Name(),
		PassportPhoto:              temp5.Name(),
		BusinessPremiseCertificate: temp6.Name(),
		FormC07:                    temp7.Name(),
		IDDocument:                 temp8.Name(),
		Doctype:                    "organization",
	}

	jsonData, err := json.Marshal(org)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(jUser))
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
