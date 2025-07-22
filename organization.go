package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// delete member
func reportMemberHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing member ID", http.StatusBadRequest)
		return
	}

	// Fetch organization
	orgResp, err := findOrganization(dbFindUrl, id)
	if err != nil || len(orgResp.Body) == 0 {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	org := orgResp.Body[0]
	org.Status = "Reported"

	jsonData, err := json.Marshal(org)
	if err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", dbUrl, org.ID), bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create update request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin_members", http.StatusSeeOther)
}

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
	tmpl.ExecuteTemplate(w, "organization_org.html", organization)
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

func org(w http.ResponseWriter, r *http.Request) {
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
	organization.CompanyLogo = strings.ReplaceAll(organization.CompanyLogo, "\\", "/")

	tmpl.ExecuteTemplate(w, "org.html", organization)
}

func organizationRegister(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "organization_register.html", nil)
}

func organizationRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var referee1 Referee1
	var referee2 Referee2
	year, month, day := time.Now().Date()
	r.ParseMultipartForm(10 * 1024 * 1024)
	name := r.FormValue("companyName")
	email := strings.ToLower(r.FormValue("email"))
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

	referee1.ID = uuid.NewString()
	referee1.Name = r.FormValue("referee1Name")
	referee1.BusinessName = r.FormValue("referee1Business")
	referee1.PhoneNumber = r.FormValue("referee1Phone")
	referee1.ChamberRegNumber = r.FormValue("referee1RegNumber")

	referee2.ID = uuid.NewString()
	referee2.Name = r.FormValue("referee2Name")
	referee2.BusinessName = r.FormValue("referee2Business")
	referee2.PhoneNumber = r.FormValue("referee2Phone")
	referee2.ChamberRegNumber = r.FormValue("referee2RegNumber")

	orgResp, err := findOrganization(dbFindUrl, email)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(orgResp.Body) != 0 {
		tmpl.ExecuteTemplate(w, "organization_register.html", "An Organization is Already Using This Email, Please Choose Another one.")
		return
	}

	coverLetter, _, err := r.FormFile("coverLetter")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer coverLetter.Close()
	coverLetterTemp, err2 := os.CreateTemp("static/Cover_Letters", name+" Cover Letter-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer coverLetterTemp.Close()
	coverLetterBytes, err3 := io.ReadAll(coverLetter)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	coverLetterTemp.Write(coverLetterBytes)

	memorandum, _, err := r.FormFile("memorandum")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer memorandum.Close()
	memorandumTemp, err2 := os.CreateTemp("static/Memorandums", name+" Memorandum-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer memorandumTemp.Close()
	memorandumBytes, err3 := io.ReadAll(memorandum)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	memorandumTemp.Write(memorandumBytes)

	businessCertificate, _, err := r.FormFile("registrationCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer businessCertificate.Close()
	businessCertificateTemp, err2 := os.CreateTemp("static/Business_Name_Certs", name+" Business Name Cert-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer businessCertificateTemp.Close()
	businessCertificateBytes, err3 := io.ReadAll(businessCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	businessCertificateTemp.Write(businessCertificateBytes)

	incorporationCertificate, _, err := r.FormFile("incorporationCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer incorporationCertificate.Close()
	incorporationCertTemp, err2 := os.CreateTemp("static/Incorporation_Certs", name+" Incorporation Cert-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer incorporationCertTemp.Close()
	incorporationCertBytes, err3 := io.ReadAll(incorporationCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	incorporationCertTemp.Write(incorporationCertBytes)

	companyLogo, _, err := r.FormFile("companyLogo")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer companyLogo.Close()
	companyLogoTemp, err2 := os.CreateTemp("static/Company_Logos", name+" Company Logo-*.png")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer companyLogoTemp.Close()
	companyLogoBytes, err3 := io.ReadAll(companyLogo)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	companyLogoTemp.Write(companyLogoBytes)

	businessPremiseCertificate, _, err := r.FormFile("premisesCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer businessPremiseCertificate.Close()
	businessPremiseCertificateTemp, err2 := os.CreateTemp("static/Premise_Certs", name+" Premise Cert-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer businessPremiseCertificateTemp.Close()
	businessPremiseCertificateBytes, err3 := io.ReadAll(businessPremiseCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	businessPremiseCertificateTemp.Write(businessPremiseCertificateBytes)

	formC07, _, err := r.FormFile("formC07")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer formC07.Close()
	formC07Temp, err2 := os.CreateTemp("static/FormC07s", name+" FormC07-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer formC07Temp.Close()
	formC07Bytes, err3 := io.ReadAll(formC07)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	formC07Temp.Write(formC07Bytes)

	// idDocument := r.FormValue("idType")
	idDocument, _, err := r.FormFile("idDocument")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer idDocument.Close()
	idDocumentTemp, err2 := os.CreateTemp("static/ID_Documents", name+" ID_Document-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer idDocumentTemp.Close()
	idDocumentBytes, err3 := io.ReadAll(idDocument)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	idDocumentTemp.Write(idDocumentBytes)
	idType := r.FormValue("idType")

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
		CoverLetter:                        coverLetterTemp.Name(),
		CoverLetterApproval:                "Pending",
		Memorandum:                         memorandumTemp.Name(),
		MemorandumApproval:                 "Pending",
		BusinessCertificate:                businessCertificateTemp.Name(),
		BusinessCertificateApproval:        "Pending",
		IncorporationCertificate:           incorporationCertTemp.Name(),
		IncorporationCertificateApproval:   "Pending",
		CompanyLogo:                        companyLogoTemp.Name(),
		CompanyLogoApproval:                "Pending",
		BusinessPremiseCertificate:         businessPremiseCertificateTemp.Name(),
		BusinessPremiseCertificateApproval: "Pending",
		FormC07:                            formC07Temp.Name(),
		FormC07Approval:                    "Pending",
		IDDocument:                         idDocumentTemp.Name(),
		IDDocumentApproval:                 "Pending",
		IDDocumentType:                     idType,
		Year:                               year,
		Month:                              month.String(),
		Day:                                day,
		DateJoined:                         dateJoined.String(),
		UserId:                             currentUser.ID,
		Referee1:                           referee1,
		Referee2:                           referee2,
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
	organizationName := r.FormValue("organizationName")
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
			// fmt.Println("Already Paid For Current Year")
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
	temp, err2 := os.CreateTemp("static/Payments", organizationName+" payment-*.pdf")
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

func updateRejectedDocuments(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "update_rejected_documents.html", organization)
}

func updateMemorandum(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "update_memorandum.html", nil)
}

func updateMemorandumHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "404.html", "Organization Not Found")
		return
	}
	organization := orgResp.Body[0]
	memorandum, _, err := r.FormFile("memorandum")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer memorandum.Close()
	memorandumTemp, err2 := os.CreateTemp("static/Memorandums", organization.Name+" Memorandum-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer memorandumTemp.Close()
	memorandumBytes, err3 := io.ReadAll(memorandum)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	memorandumTemp.Write(memorandumBytes)

	organization.Memorandum = memorandumTemp.Name()

	jsonData, err := json.Marshal(&organization)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+organization.ID, bytes.NewBuffer(jsonData))
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
	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func updateCoverLetter(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "update_cover_letter.html", nil)
}

func updateCoverLetterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "400.html", "Organization Not Found")
		return
	}
	organization := orgResp.Body[0]
	coverLetter, _, err := r.FormFile("coverLetter")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer coverLetter.Close()
	coverLetterTemp, err2 := os.CreateTemp("static/Cover_Letters", organization.Name+" Cover Letter-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer coverLetterTemp.Close()
	coverLetterBytes, err3 := io.ReadAll(coverLetter)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	coverLetterTemp.Write(coverLetterBytes)

	organization.CoverLetter = coverLetterTemp.Name()

	jsonData, err := json.Marshal(&organization)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+organization.ID, bytes.NewBuffer(jsonData))
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
	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func updateBusinessCertificate(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "update_business_certificate.html", nil)
}

func updateBusinessCertificateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "404.html", "Organization Not Found")
		return
	}
	organization := orgResp.Body[0]
	businessCertificate, _, err := r.FormFile("registrationCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer businessCertificate.Close()
	businessCertificateTemp, err2 := os.CreateTemp("static/Business_Name_Certs", organization.Name+" Business Name Cert-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer businessCertificateTemp.Close()
	businessCertificateBytes, err3 := io.ReadAll(businessCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	businessCertificateTemp.Write(businessCertificateBytes)
	organization.BusinessCertificate = businessCertificateTemp.Name()
	jsonData, err := json.Marshal(&organization)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+organization.ID, bytes.NewBuffer(jsonData))
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
	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func updateIncorporationCertificate(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "update_incorporation_certificate.html", nil)
}

func updateIncorporationCertificateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "404.html", "Organization Not Found")
		return
	}
	organization := orgResp.Body[0]
	incorporationCertificate, _, err := r.FormFile("incorporationCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer incorporationCertificate.Close()
	incorporationCertTemp, err2 := os.CreateTemp("static/Incorporation_Certs", organization.Name+" Incorporation Cert-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer incorporationCertTemp.Close()
	incorporationCertBytes, err3 := io.ReadAll(incorporationCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	incorporationCertTemp.Write(incorporationCertBytes)

	organization.IncorporationCertificate = incorporationCertTemp.Name()
	jsonData, err := json.Marshal(&organization)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+organization.ID, bytes.NewBuffer(jsonData))
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
	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func updateBusinessPremiseCertificate(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "update_business_premise_certificate.html", nil)
}

func updateBusinessPremiseCertificateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "404.html", "Organization Not Found")
		return
	}
	organization := orgResp.Body[0]
	businessPremiseCertificate, _, err := r.FormFile("premisesCert")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer businessPremiseCertificate.Close()
	businessPremiseCertificateTemp, err2 := os.CreateTemp("static/Premise_Certs", organization.Name+" Premise Cert-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer businessPremiseCertificateTemp.Close()
	businessPremiseCertificateBytes, err3 := io.ReadAll(businessPremiseCertificate)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	businessPremiseCertificateTemp.Write(businessPremiseCertificateBytes)

	organization.BusinessPremiseCertificate = businessPremiseCertificateTemp.Name()
	jsonData, err := json.Marshal(&organization)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+organization.ID, bytes.NewBuffer(jsonData))
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
	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func updateCompanyLogo(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "update_company_logo.html", nil)
}

func updateCompanyLogoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "404.html", "Organization Not Found")
		return
	}
	organization := orgResp.Body[0]
	companyLogo, _, err := r.FormFile("companyLogo")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer companyLogo.Close()
	companyLogoTemp, err2 := os.CreateTemp("static/Company_Logos", organization.Name+" Company Logo-*.png")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer companyLogoTemp.Close()
	companyLogoBytes, err3 := io.ReadAll(companyLogo)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	companyLogoTemp.Write(companyLogoBytes)

	organization.CompanyLogo = companyLogoTemp.Name()
	jsonData, err := json.Marshal(&organization)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+organization.ID, bytes.NewBuffer(jsonData))
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
	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func updateFormC07(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "update_formco7.html", nil)
}

func updateFormC07Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "404.html", "Organization Not Found")
		return
	}
	organization := orgResp.Body[0]
	formC07, _, err := r.FormFile("formC07")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer formC07.Close()
	formC07Temp, err2 := os.CreateTemp("static/FormC07s", organization.Name+" FormC07-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer formC07Temp.Close()
	formC07Bytes, err3 := io.ReadAll(formC07)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	formC07Temp.Write(formC07Bytes)

	organization.FormC07 = formC07Temp.Name()
	jsonData, err := json.Marshal(&organization)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+organization.ID, bytes.NewBuffer(jsonData))
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
	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func updateIDDocument(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "update_idDocument.html", nil)
}

func updateIDDocumentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	orgResp, err := findOrganization(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgResp.Body) == 0 {
		tmpl.ExecuteTemplate(w, "404.html", "Organization Not Found")
		return
	}
	organization := orgResp.Body[0]
	idDocument, _, err := r.FormFile("idType")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer idDocument.Close()
	idDocumentTemp, err2 := os.CreateTemp("static/ID_Documents", organization.Name+" ID_Document-*.pdf")
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	defer idDocumentTemp.Close()
	idDocumentBytes, err3 := io.ReadAll(idDocument)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	idDocumentTemp.Write(idDocumentBytes)

	organization.IDDocument = idDocumentTemp.Name()
	jsonData, err := json.Marshal(&organization)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+organization.ID, bytes.NewBuffer(jsonData))
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
	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
