package main

import (
	"fmt"
	"net/http"
)

func admin(w http.ResponseWriter, r *http.Request) {
	var totalMembers int
	type PageResult struct {
		UserCount       int
		PendingOrgCount int
		PendingOrg      []Organization
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

	orgStatusResponse, err := findOrganizationApprovalStatus(dbFindUrl, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(userResponse.Body) == 0 {
		fmt.Println("No Record Found")
		return
	}
	pendingOrg := len(orgStatusResponse.Body)
	p := PageResult{
		UserCount:       totalMembers,
		PendingOrgCount: pendingOrg,
		PendingOrg:      orgStatusResponse.Body,
	}
	tmpl.ExecuteTemplate(w, "admin.html", p)
}
