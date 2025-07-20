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
	ID                                 string   `json:"_id"`
	Name                               string   `json:"name"`
	Address                            string   `json:"address"`
	NatureOfBusiness                   string   `json:"businessNature"`
	PhoneNumber                        string   `json:"phoneNumber"`
	CACNumber                          string   `json:"cacNumber"`
	Bankers                            string   `json:"bankers"`
	NumberOfEmployees                  string   `json:"numberOfEmployees"`
	NonNigerianEmployees               string   `json:"nonNigerianEmployees"`
	NumberOfDirectors                  string   `json:"numberOfDirectors"`
	NonNigerianDirectors               string   `json:"nonNigerianDiretors"`
	ContactPerson                      string   `json:"contactPerson"`
	Representative                     string   `json:"representative"`
	Email                              string   `json:"email"`
	CoverLetter                        string   `json:"coverLetter"`
	CoverLetterApproval                string   `json:"coverLetterApproval"`
	Memorandum                         string   `json:"memorandum"`
	MemorandumApproval                 string   `json:"memorandumApproval"`
	BusinessCertificate                string   `json:"businessCertificate"`
	BusinessCertificateApproval        string   `json:"businessCertificateApproval"`
	IncorporationCertificate           string   `json:"incorporationCertificate"`
	IncorporationCertificateApproval   string   `json:"incorporationCertificateApproval"`
	BusinessPremiseCertificate         string   `json:"businessPremiseCertificate"`
	BusinessPremiseCertificateApproval string   `json:"businessPremiseCertificateApproval"`
	CompanyLogo                        string   `json:"companyLogo"`
	CompanyLogoApproval                string   `json:"companyLogoApproval"`
	FormC07                            string   `json:"formC07"`
	FormC07Approval                    string   `json:"formC07Approval"`
	IDDocument                         string   `json:"idDocument"`
	IDDocumentApproval                 string   `json:"idDocumentApproval"`
	IDDocumentType                     string   `json:"IdDocumentType"`
	Year                               int      `json:"year"`
	Month                              string   `json:"month"`
	Day                                int      `json:"day"`
	DateJoined                         string   `json:"dateJoined"`
	Status                             string   `json:"status"`
	UserId                             string   `json:"userId"`
	Referee1                           Referee1 `json:"referee1"`
	Referee2                           Referee2 `json:"referee2"`
	Rev                                string   `json:"_rev,omitempty"`
	Doctype                            string   `json:"doctype"`
}

type Referee1 struct {
	ID               string `json:"_id"`
	Name             string `json:"name"`
	PhoneNumber      string `json:"phoneNumber"`
	BusinessName     string `json:"businessName"`
	ChamberRegNumber string `json:"chamberRegNumber"`
}

type RefereeResponse struct {
	Status   string     `json:"status"`
	Body     []Referee1 `json:"body"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type Referee2 struct {
	ID               string `json:"_id"`
	Name             string `json:"name"`
	PhoneNumber      string `json:"phoneNumber"`
	BusinessName     string `json:"businessName"`
	ChamberRegNumber string `json:"chamberRegNumber"`
}

type Referee2Response struct {
	Status   string     `json:"status"`
	Body     []Referee2 `json:"body"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type OrganizationResponse struct {
	Status   string
	Body     []Organization `json:"docs"`
	Bookmark string         `json:"bookmark"`
	Warning  string         `json:"warning"`
}

type BankTransfer struct {
	ID               string `json:"_id"`
	UserId           string `json:"userId"`
	OrganizationName string `json:"name"`
	PaymentMethod    string `json:"paymentMethod"`
	Status           string `json:"status"`
	Year             int    `json:"year"`
	Month            string `json:"month"`
	Day              int    `json:"day"`
	Date             string `json:"date"`
	OrganizationId   string `json:"organizationId"`
	RecieptFile      string `json:"receiptFile"`
	Doctype          string `json:"doctype"`
	Rev              string `json:"_rev,omitempty"`
}

type BankTransferResponse struct {
	Status   string
	Body     []BankTransfer `json:"docs"`
	Bookmark string         `json:"bookmark"`
	Warning  string         `json:"warning"`
}
