package main

type User struct {
	ID       string `json:"_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Rev      string `json:"_rev,omitempty"`
	Doctype  string `json:"doctype"`
}

// type UserProfile struct {
// 	ID        string
// 	FirstName string
// 	LastName  string
// 	Image     string
// }

// type UserProfileResponse struct {
// 	Status   string
// 	Body     []UserProfile `json:"docs"`
// 	Bookmark string        `json:"bookmark"`
// 	Warning  string        `json:"warning"`
// }

type AdminProfile struct {
}

type SuperAdminProfile struct {
}

type Organization struct {
	ID                         string `json:"_id"`
	Name                       string `json:"name"`
	Address                    string `json:"address"`
	NatureOfBusiness           string `json:"businessNature"`
	Bankers                    string `json:"bankers"`
	NumberOfEmployees          string `json:"numberOfEmployees"`
	NonNigerianEmployees       string `json:"nonNigerianEmployees"`
	NumberOfDirectors          string `json:"numberOfDirectors"`
	NonNigerianDirectors       string `json:"nonNigerianDiretors"`
	ContactPerson              string `json:"contactPerson"`
	Representative             string `json:"representative"`
	Email                      string `json:"email"`
	CoverLetter                string `json:"coverLetter"`
	Memorandum                 string `json:"memorandum"`
	BusinessCertificate        string `json:"businessCertificate"`
	IncorporationCertificate   string `json:"incorporationCertificate"`
	BusinessPremiseCertificate string `json:"businessPremiseCertificate"`
	PassportPhoto              string `json:"passportPhoto"`
	FormC07                    string `json:"formC07"`
	IDDocument                 string `json:"idDocument"`
	DateJoined                 string `json:"dateJoined"`
	Approved                   string `json:"approved"`
	UserId                     string `json:"userId"`
	Doctype                    string `json:"doctype"`
}

type OrganizationResponse struct {
	Status   string
	Body     []Organization `json:"docs"`
	Bookmark string         `json:"bookmark"`
	Warning  string         `json:"warning"`
}
