package main

import (
	"fmt"
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
	tmpl.ExecuteTemplate(w, "admin_documents.html", nil)
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
	fmt.Println(organizations)
	tmpl.ExecuteTemplate(w, "admin_members.html", organizations)
}

func adminPayment(w http.ResponseWriter, r *http.Request) {
	var pendingBankTransfer []BankTransfer
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
	tmpl.ExecuteTemplate(w, "admin_payments.html", pendingBankTransfer)
}

func adminReport(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin_report.html", nil)
}

func adminSettings(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin_settings.html", nil)
}
