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

func findUsers(url string) (UserResponse, error) {
	var userResponse UserResponse
	jsonData := map[string]map[string]any{"selector": {"doctype": "user", "role": "user"}}
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

// func findUserProfile(url string, userId string) (UserProfileResponse, error) {
// 	var response UserProfileResponse
// 	jsonData := map[string]map[string]any{"selector": {"_id": userId, "doctype": "profile"}}
// 	data, error := json.Marshal(jsonData)
// 	if error != nil {
// 		log.Fatalln("Marshal", error)
// 	}
// 	request, error := http.NewRequest("POST", url, bytes.NewBuffer(data))
// 	if error != nil {
// 		fmt.Println("Byte Error", error)
// 	}
// 	request.Header.Set("Content-type", "application/json")
// 	client := &http.Client{}
// 	res, error := client.Do(request)
// 	if error != nil {
// 		fmt.Println("Req Err", error)
// 	}
// 	defer res.Body.Close()
// 	body, _ := io.ReadAll(res.Body)
// 	error = json.Unmarshal(body, &response)
// 	if error != nil {
// 		log.Fatalln("UnMarshal Err: ", error)
// 	}
// 	return response, error
// }

func findAllOrganizations(url string) (OrganizationResponse, error) {
	var response OrganizationResponse
	jsonData := map[string]map[string]any{"selector": {"doctype": "organization"}}
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

func findOrganization(url, userId string) (OrganizationResponse, error) {
	var response OrganizationResponse
	jsonData := map[string]map[string]any{"selector": {"userId": userId, "doctype": "organization"}}
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

func findOrganizationById(url, id string) (OrganizationResponse, error) {
	var response OrganizationResponse
	jsonData := map[string]map[string]any{"selector": {"_id": id, "doctype": "organization"}}
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

func findOrganizationApprovalStatus(url, status string) (OrganizationResponse, error) {
	var response OrganizationResponse
	jsonData := map[string]map[string]any{"selector": {"status": status, "doctype": "organization"}}
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

func findBankPaymentApprovalStatus(url, status string) (BankTransferResponse, error) {
	var response BankTransferResponse
	jsonData := map[string]map[string]any{"selector": {"status": status, "doctype": "bankTransfer"}}
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

func findOrganizationPayment(url, organizationId string) (BankTransferResponse, error) {
	var response BankTransferResponse
	jsonData := map[string]map[string]any{"selector": {"organizationId": organizationId, "doctype": "bankTransfer"}}
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

func findOrganizationPaymentByID(url, paymentId string) (BankTransferResponse, error) {
	var response BankTransferResponse
	jsonData := map[string]map[string]any{"selector": {"_id": paymentId, "doctype": "bankTransfer"}}
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

func findOrganizationPayments(url string) (BankTransferResponse, error) {
	var response BankTransferResponse
	jsonData := map[string]map[string]any{"selector": {"doctype": "bankTransfer"}}
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

func findOrganizationReports(url, id string) (ReportResponse, error) {
	var response ReportResponse
	jsonData := map[string]map[string]any{"selector": {"organizationId": id, "doctype": "report"}}
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
