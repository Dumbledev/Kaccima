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
	fmt.Println(bankTransferPendingStatusResponse.Body, "Bodus")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(bankTransferPendingStatusResponse.Body) == 0 {
		fmt.Println("No Record Found")
		return
	}
	pendingBankTransfer := bankTransferPendingStatusResponse.Body
	p := PageResult{
		UserCount:          totalMembers,
		PendingOrgCount:    pendingOrg,
		PendingOrg:         orgPendingStatusResponse.Body,
		PendingBankPayment: pendingBankTransfer,
	}
	tmpl.ExecuteTemplate(w, "admin.html", p)
}
