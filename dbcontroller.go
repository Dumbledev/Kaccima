package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type UserResponse struct {
	Status   string
	Body     []User `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
}

func findUser(url, email string) (UserResponse, error) {
	var userResponse UserResponse
	jsonData := map[string]map[string]any{"selector": {"email": email, "doctype": "user"}}
	data, error := json.Marshal(jsonData)
	if error != nil {
		log.Fatalln("Marshal", error)
	}
	request, error := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if error != nil {
		fmt.Println("Byte Error", error)
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		fmt.Println("Req Err", error)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	error = json.Unmarshal(body, &userResponse)
	if error != nil {
		log.Fatalln("UnMarshal Err: ", error)
	}
	return userResponse, error
}

func findUserByEmail(url string, email string) (UserResponse, error) {
	var userResponse UserResponse
	jsonData := map[string]map[string]any{"selector": {"email": email, "doctype": "user"}}
	data, error := json.Marshal(jsonData)
	if error != nil {
		log.Fatalln("Marshal", error)
	}
	request, error := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if error != nil {
		fmt.Println("Byte Error", error)
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		fmt.Println("Req Err", error)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	error = json.Unmarshal(body, &userResponse)
	if error != nil {
		log.Fatalln("UnMarshal Err: ", error)
	}
	return userResponse, error
}

func findUserById(url string, userId string) (UserResponse, error) {
	var userResponse UserResponse
	jsonData := map[string]map[string]any{"selector": {"_id": userId, "doctype": "user"}}
	data, error := json.Marshal(jsonData)
	if error != nil {
		log.Fatalln("Marshal", error)
	}
	request, error := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if error != nil {
		fmt.Println("Byte Error", error)
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		fmt.Println("Req Err", error)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	error = json.Unmarshal(body, &userResponse)
	if error != nil {
		log.Fatalln("UnMarshal Err: ", error)
	}
	return userResponse, error
}

func findUserProfile(url string, userId string) (UserProfileResponse, error) {
	var response UserProfileResponse
	jsonData := map[string]map[string]any{"selector": {"_id": userId, "doctype": "profile"}}
	data, error := json.Marshal(jsonData)
	if error != nil {
		log.Fatalln("Marshal", error)
	}
	request, error := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if error != nil {
		fmt.Println("Byte Error", error)
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		fmt.Println("Req Err", error)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	error = json.Unmarshal(body, &response)
	if error != nil {
		log.Fatalln("UnMarshal Err: ", error)
	}
	return response, error
}

func findOrganization(url, email string) (OrganizationResponse, error) {
	var response OrganizationResponse
	jsonData := map[string]map[string]any{"selector": {"email": email, "doctype": "organization"}}
	data, error := json.Marshal(jsonData)
	if error != nil {
		log.Fatalln("Marshal", error)
	}
	request, error := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if error != nil {
		fmt.Println("Byte Error", error)
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		fmt.Println("Req Err", error)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	error = json.Unmarshal(body, &response)
	if error != nil {
		log.Fatalln("UnMarshal Err: ", error)
	}
	return response, error
}
