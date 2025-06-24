package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func admin(w http.ResponseWriter, r *http.Request) {
	var totalMembers int
	type PageResult struct {
		UserCount          int
		PendingOrgCount    int
		PendingOrg         []Organization
		PendingBankPayment []BankTransfer
	}
	var pendingBankTransfer []BankTransfer
	userResponse, err := findUsers(dbFindUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(userResponse.Body) == 0 {
		fmt.Println("No Record Found")
		return
	}
	totalMembers = len(userResponse.Body)

	orgPendingStatusResponse, err := findOrganizationApprovalStatus(dbFindUrl, "Pending")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(userResponse.Body) == 0 {
		fmt.Println("No Record Found")
		return
	}
	pendingOrg := len(orgPendingStatusResponse.Body)

	bankTransferPendingStatusResponse, err := findBankPaymentApprovalStatus(dbFindUrl, "Pending")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(bankTransferPendingStatusResponse.Body) == 0 {
		pendingBankTransfer = []BankTransfer{}
		fmt.Println("No Record Found")
	} else {
		pendingBankTransfer = bankTransferPendingStatusResponse.Body
	}

	p := PageResult{
		UserCount:          totalMembers,
		PendingOrgCount:    pendingOrg,
		PendingOrg:         orgPendingStatusResponse.Body,
		PendingBankPayment: pendingBankTransfer,
	}
	tmpl.ExecuteTemplate(w, "admin.html", p)
}

func adminDocuments(w http.ResponseWriter, r *http.Request) {
	organizations := []Organization{}
	orgResp, err := findAllOrganizations(dbFindUrl)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organizations = orgResp.Body
	}
	tmpl.ExecuteTemplate(w, "admin_documents.html", organizations)
}

func adminOrganizationDocuments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, id)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	tmpl.ExecuteTemplate(w, "admin_organization_documents.html", organization)
}

func adminMembers(w http.ResponseWriter, r *http.Request) {
	organizations := []Organization{}
	orgResp, err := findAllOrganizations(dbFindUrl)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organizations = orgResp.Body
	}
	tmpl.ExecuteTemplate(w, "admin_members.html", organizations)
}

func adminPayment(w http.ResponseWriter, r *http.Request) {
	var bankTransfers []BankTransfer
	bankTransferPendingStatusResponse, err := findOrganizationPayments(dbFindUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(bankTransferPendingStatusResponse.Body) == 0 {
		bankTransfers = []BankTransfer{}
		fmt.Println("No Record Found")
	} else {
		bankTransfers = bankTransferPendingStatusResponse.Body
	}
	tmpl.ExecuteTemplate(w, "admin_payments.html", bankTransfers)
}

func adminReport(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin_report.html", nil)
}

func adminSettings(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin_settings.html", nil)
}

func acceptReceipt(w http.ResponseWriter, r *http.Request) {
	paymentId := r.PathValue("paymentId")
	fmt.Println(paymentId, "id")
	payment := BankTransfer{}
	paymentResp, err := findOrganizationPaymentByID(dbFindUrl, paymentId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(paymentResp.Body) != 0 {
		payment = paymentResp.Body[0]
	}
	payment.Status = "Accepted"
	fmt.Println(payment)

	jsonData, err := json.Marshal(&payment)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+payment.ID, bytes.NewBuffer(jsonData))
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
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
	http.Redirect(w, r, "/admin", http.StatusPermanentRedirect)
}

func rejectReceipt(w http.ResponseWriter, r *http.Request) {
	paymentId := r.PathValue("paymentId")
	payment := BankTransfer{}
	paymentResp, err := findOrganizationPaymentByID(dbFindUrl, paymentId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(paymentResp.Body) != 0 {
		payment = paymentResp.Body[0]
	}
	payment.Status = "Rejected"

	jsonData, err := json.Marshal(&payment)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+payment.ID, bytes.NewBuffer(jsonData))
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
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
	http.Redirect(w, r, "/admin", http.StatusPermanentRedirect)
}
