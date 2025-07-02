package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func organizationDashboard(w http.ResponseWriter, r *http.Request) {
	type PageResult struct {
		Organization Organization
		Payment      BankTransfer
	}
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) == 0 {
		http.Redirect(w, r, "/organization_register", http.StatusPermanentRedirect)
		return
	}
	organization := orgResp.Body[0]
	paymentResp, err := findOrganizationPayment(dbFindUrl, organization.ID)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) == 0 {
		http.Redirect(w, r, "/organization_register", http.StatusPermanentRedirect)
		return
	}
	payment := paymentResp.Body[0]
	p := PageResult{
		Organization: organization,
		Payment:      payment,
	}
	tmpl.ExecuteTemplate(w, "dashboard.html", p)
}

func organization(w http.ResponseWriter, r *http.Request) {
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) == 0 {
		http.Redirect(w, r, "/organization_register", http.StatusPermanentRedirect)
		return
	}
	organization := orgResp.Body[0]
	tmpl.ExecuteTemplate(w, "organization_profile.html", organization)
}

func notification(w http.ResponseWriter, r *http.Request) {
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) == 0 {
		http.Redirect(w, r, "/organization_register", http.StatusPermanentRedirect)
		return
	}
	organization := orgResp.Body[0]
	tmpl.ExecuteTemplate(w, "notification.html", organization)
}

func reviewedDocuments(w http.ResponseWriter, r *http.Request) {
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) == 0 {
		http.Redirect(w, r, "/organization_register", http.StatusPermanentRedirect)
		return
	}
	organization := orgResp.Body[0]
	tmpl.ExecuteTemplate(w, "reviewedDocuments.html", organization)
}

func profile(w http.ResponseWriter, r *http.Request) {
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) == 0 {
		http.Redirect(w, r, "/organization_register", http.StatusPermanentRedirect)
		return
	}
	organization := orgResp.Body[0]
	tmpl.ExecuteTemplate(w, "profile.html", organization)
}

func organizationRegister(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "organization_register.html", nil)
}

func organizationRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicaton/json")
	year, month, day := time.Now().Date()
	r.ParseMultipartForm(10 * 1024 * 1024)
	name := r.FormValue("companyName")
	email := r.FormValue("email")
	address := r.FormValue("officeAddress")
	employeesNo := r.FormValue("employeesNo")
	nonNigerianEmployees := r.FormValue("nonNigerianEmployees")
	directorsNo := r.FormValue("directorsNo")
	nonNigerianDirectors := r.FormValue("nonNigerianDirectors")
	businessNature := r.FormValue("businessNature")
	bankers := r.FormValue("bankers")
	contactPerson := r.FormValue("contactPerson")
	rep := r.FormValue("representative")
	cacNumber := r.FormValue("cacNumber")
	phoneNumber := r.FormValue("phoneNumber")

	dateJoined := time.Now().Local()

	orgResp, err := findOrganization(dbFindUrl, email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(orgResp.Body) != 0 {
		tmpl.ExecuteTemplate(w, "organization_register.html", "Organization Already Registered, Please Choose Another one.")
		return
	}

	coverLetter, _, err := r.FormFile("coverLetter")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer coverLetter.Close()
	temp1, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp1.Close()
	fileBytes, err3 := io.ReadAll(coverLetter)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp1.Write(fileBytes)

	memorandum, _, err := r.FormFile("memorandum")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer memorandum.Close()
	temp2, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp2.Close()
	fileBytes1, err3 := io.ReadAll(memorandum)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp2.Write(fileBytes1)

	businessCertificate, _, err := r.FormFile("registrationCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer businessCertificate.Close()
	temp3, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp3.Close()
	fileBytes2, err3 := io.ReadAll(businessCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp2.Write(fileBytes2)

	incorporationCertificate, _, err := r.FormFile("incorporationCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer incorporationCertificate.Close()
	temp4, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp4.Close()
	fileBytes3, err3 := io.ReadAll(incorporationCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp2.Write(fileBytes3)

	passportPhoto, _, err := r.FormFile("passportPhotos")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer passportPhoto.Close()
	temp5, err2 := os.CreateTemp("static", "file-*.png")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp5.Close()
	fileBytes4, err3 := io.ReadAll(passportPhoto)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp5.Write(fileBytes4)

	businessPremiseCertificate, _, err := r.FormFile("premisesCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer businessPremiseCertificate.Close()
	temp6, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp6.Close()
	fileBytes5, err3 := io.ReadAll(businessPremiseCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp6.Write(fileBytes5)

	formC07, _, err := r.FormFile("formC07")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer formC07.Close()
	temp7, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp7.Close()
	fileBytes6, err3 := io.ReadAll(formC07)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp7.Write(fileBytes6)

	idDocument, _, err := r.FormFile("nationalId")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer idDocument.Close()
	temp8, err2 := os.CreateTemp("static", "file-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp8.Close()
	fileBytes7, err3 := io.ReadAll(idDocument)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp5.Write(fileBytes7)

	var org = Organization{
		ID:                                 uuid.NewString(),
		Name:                               name,
		Email:                              email,
		Address:                            address,
		NumberOfEmployees:                  employeesNo,
		NonNigerianEmployees:               nonNigerianEmployees,
		NumberOfDirectors:                  directorsNo,
		NonNigerianDirectors:               nonNigerianDirectors,
		NatureOfBusiness:                   businessNature,
		Bankers:                            bankers,
		ContactPerson:                      contactPerson,
		Representative:                     rep,
		CACNumber:                          cacNumber,
		PhoneNumber:                        phoneNumber,
		CoverLetter:                        temp1.Name(),
		CoverLetterApproval:                "Pending",
		Memorandum:                         temp2.Name(),
		MemorandumApproval:                 "Pending",
		BusinessCertificate:                temp3.Name(),
		BusinessCertificateApproval:        "Pending",
		IncorporationCertificate:           temp4.Name(),
		IncorporationCertificateApproval:   "Pending",
		PassportPhoto:                      temp5.Name(),
		PassportPhotoApproval:              "Pending",
		BusinessPremiseCertificate:         temp6.Name(),
		BusinessPremiseCertificateApproval: "Pending",
		FormC07:                            temp7.Name(),
		FormC07Approval:                    "Pending",
		IDDocument:                         temp8.Name(),
		IDDocumentApproval:                 "Pending",
		Year:                               year,
		Month:                              month.String(),
		Day:                                day,
		DateJoined:                         dateJoined.String(),
		UserId:                             currentUser.ID,
		Status:                             "Pending",
		Doctype:                            "organization",
	}

	jsonData, err := json.Marshal(org)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Server Error(1)")
		return
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		tmpl.ExecuteTemplate(w, "sign_up.html", "Server Error(2)")
		return
	}
	defer res.Body.Close()

	http.Redirect(w, r, "/payment", http.StatusSeeOther)
}

func payment(w http.ResponseWriter, r *http.Request) {
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) == 0 {
		http.Redirect(w, r, "/organization_register", http.StatusPermanentRedirect)
		return
	}
	organization := orgResp.Body[0]
	tmpl.ExecuteTemplate(w, "payment.html", organization)
}

func bankTransferHandler(w http.ResponseWriter, r *http.Request) {
	payment := BankTransfer{}
	year, month, day := time.Now().Date()
	r.ParseMultipartForm(10 * 1024 * 1024)
	organizationId := r.FormValue("organizationId")
	name := r.FormValue("name")

	paymentResp, err := findOrganizationPayment(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(paymentResp.Body) != 0 {
		payment = paymentResp.Body[0]
		if payment.Year == year {
			http.Redirect(w, r, "/payment", http.StatusPermanentRedirect)
			fmt.Println("Already Paid For Current Year")
			return
		}
	}
	// payment := paymentResp.Body[0]

	reciept, _, err := r.FormFile("receiptFile")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer reciept.Close()
	temp, err2 := os.CreateTemp("static", "payment-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer temp.Close()
	fileBytes, err3 := io.ReadAll(reciept)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	temp.Write(fileBytes)
	bankTranfer := BankTransfer{
		ID:               uuid.NewString(),
		UserId:           currentUser.ID,
		OrganizationName: name,
		OrganizationId:   organizationId,
		PaymentMethod:    "Bank Transfer",
		Status:           "Pending",
		Date:             time.Now().String(),
		Day:              day,
		Month:            month.String(),
		Year:             year,
		Doctype:          "bankTransfer",
		RecieptFile:      temp.Name(),
	}
	jsonData, err := json.Marshal(bankTranfer)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Server Error(1)")
		return
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		tmpl.ExecuteTemplate(w, "sign_up.html", "Server Error(2)")
		return
	}
	defer res.Body.Close()
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
