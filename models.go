package main

type User struct {
	ID       string `json:"_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Approved bool   `json:"approved"`
	Rev      string `json:"_rev,omitempty"`
	Doctype  string `json:"doctype"`
}

type UserProfile struct {
	ID        string
	FirstName string
	LastName  string
	Image     string
}

type AdminProfile struct {
}

type SuperAdminProfile struct {
}
