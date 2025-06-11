package main

import (
	"fmt"
	"net/http"
)

func admin(w http.ResponseWriter, r *http.Request) {
	var totalMembers int

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
	tmpl.ExecuteTemplate(w, "admin.html", totalMembers)
}
