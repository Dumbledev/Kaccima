package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var dbUrl = "http://admin:admin@localhost:5984/kaccima"
var dbFindUrl = dbUrl + "/_find"
var jwtSecret = []byte("mysecretkey")
var currentUser User

func main() {
	mux := http.NewServeMux()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
	handler := cors.Handler(mux)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/api/v1/register", register)
	mux.HandleFunc("/api/v1/login", login)

	fmt.Println("Server Running on Port :8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userResponse, err := findUserByEmail(dbFindUrl, user.Email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println(userResponse.Body)
	if len(userResponse.Body) != 0 {
		fmt.Println("Email Already In Use, Please Choose Another one.")
		return
	}
	hashedPassword, hashedErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashedErr != nil {
		fmt.Println(hashedErr)
		return
	}

	userReg := User{
		ID:       uuid.NewString(),
		Email:    user.Email,
		Password: string(hashedPassword),
		Role:     user.Role,
		Approved: false,
		Doctype:  "user",
	}

	jUser, err := json.Marshal(userReg)
	if err != nil {
		fmt.Println(err)
		return
	}
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
	userResponse, err = findUserByEmail(dbFindUrl, user.Email)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(userResponse.Body) == 0 {
		fmt.Println("No User Data Found")
		return
	}
	user = userResponse.Body[0]
	json.NewEncoder(w).Encode(user)
}

func login(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Email    string
		Password string
	}
	var body Body
	var user User

	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	userResponse, err := findUserByEmail(dbFindUrl, body.Email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(userResponse.Body) == 0 {
		fmt.Println("Invalid Username or Password")
		return
	}
	user = userResponse.Body[0]
	hash_err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if hash_err != nil {
		fmt.Println(hash_err.Error())
		http.Error(w, "Incorrect Password", http.StatusNotFound)
		return
	}

	fmt.Println("Hash Success")
	// //Generate Jwt Token7
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		http.Error(w, "Failed to create token", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(&tokenString)
}

func isAuthenticated(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, " ")
		if len(splitToken) != 2 {
			fmt.Printf("Error Reading Token")
			return
		}
		reqToken = strings.TrimSpace(splitToken[1])

		if reqToken != "" {
			// fmt.Println("Header is set! We can serve content!")

			token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if float64(time.Now().Unix()) > claims["exp"].(float64) {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				var user User
				userResponse, err := findUserById(dbFindUrl, fmt.Sprint(claims["sub"]))
				if err != nil {
					fmt.Println("err", err)
					return
				}
				if len(userResponse.Body) == 0 {
					fmt.Println("Invalid Username or Password")
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				user = userResponse.Body[0]
				currentUser = user
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			endpoint(w, r)
		} else {
			fmt.Println("Not Authorized")
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}
