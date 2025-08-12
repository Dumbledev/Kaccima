package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"html/template"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var dbUrl = "http://admin:admin@localhost:5984/kaccima"
var dbFindUrl = dbUrl + "/_find"
var store = sessions.NewCookieStore([]byte("Secret"))
var currentUser User
var tmpl *template.Template

func main() {
	tmpl, _ = template.ParseGlob("templates/*.html")

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	http.Handle("/", isAuthenticated(index))
	http.HandleFunc("/home", home)
	http.HandleFunc("/kaccimanew2025", adminSignUp)
	http.HandleFunc("/kaccimanew2025_handler", adminSignUpHandler)
	http.HandleFunc("/kaccimanew2025super", superAdminSignUp)
	http.HandleFunc("/kaccimanew2025super_handler", superAdminSignUpHandler)

	http.HandleFunc("/sign_up", signUp)
	http.HandleFunc("/sign_up_handler", signUpHandler)
	http.HandleFunc("/sign_in", signIn)
	http.HandleFunc("/sign_in_handler", signInHandler)
	http.HandleFunc("/sign_out", signOut)
	http.HandleFunc("/forgot_password", forgotPassword)
	http.HandleFunc("/forgot_password_handler", forgotPasswordHandler)
	http.HandleFunc("/password_reset_qa/{email}", passwordResetQa)
	http.HandleFunc("/password_reset_qa_handler", passwordResetQaHandler)
	http.HandleFunc("/reset_password/{email}", resetPassword)
	http.HandleFunc("/reset_password_handler", resetPasswordHandler)

	http.Handle("/dashboard", isAuthenticated(organizationDashboard))
	http.Handle("/notification", isAuthenticated(notification))
	http.Handle("/payment", isAuthenticated(payment))
	http.Handle("/bank_transfer_handler", isAuthenticated(bankTransferHandler))
	http.Handle("/update_bank_transfer", isAuthenticated(updateBankPayment))
	http.Handle("/update_bank_transfer_handler", isAuthenticated(updateBankPaymentHandler))
	http.Handle("/reviewedDocuments", isAuthenticated(reviewedDocuments))
	http.Handle("/organization", isAuthenticated(organization))
	http.Handle("/organization_register", isAuthenticated(organizationRegister))
	http.Handle("/organization_register_handler", isAuthenticated(organizationRegisterHandler))

	http.Handle("/admin", isAuthenticated(admin))
	http.Handle("/admin_documents", isAuthenticated(adminDocuments))
	http.Handle("/admin_organization_documents/{organizationId}", isAuthenticated(adminOrganizationDocuments))
	http.Handle("/admin_members", isAuthenticated(adminMembers))
	http.Handle("/admin_payments", isAuthenticated(adminPayment))
	http.Handle("/admin_report_member/{organizationId}", isAuthenticated(adminReport))
	http.Handle("/admin_report_member_handler", isAuthenticated(adminReportHandler))
	http.Handle("/admin_settings", isAuthenticated(adminSettings))

	http.Handle("/approval_admin", isAuthenticated(approvalAdmin))

	http.Handle("/accept_receipt/{paymentId}", isAuthenticated(acceptReceipt))
	http.Handle("/reject_receipt/{paymentId}", isAuthenticated(rejectReceipt))
	http.Handle("/accept_organization/{organizationId}", isAuthenticated(acceptOrganization))
	http.Handle("/reject_organization/{organizationId}", isAuthenticated(rejectOrganization))

	http.Handle("/approve_memorandum/{organizationId}", isAuthenticated(approveMemorandum))
	http.Handle("/reject_memorandum/{organizationId}", isAuthenticated(rejectMemorandum))
	http.Handle("/approve_coverLetter/{organizationId}", isAuthenticated(approveCoverLetter))
	http.Handle("/reject_coverLetter/{organizationId}", isAuthenticated(rejectCoverLetter))
	http.Handle("/approve_businessCert/{organizationId}", isAuthenticated(approveBusinessCertificate))
	http.Handle("/reject_businessCert/{organizationId}", isAuthenticated(rejectBusinessCertificate))
	http.Handle("/approve_incorporationCert/{organizationId}", isAuthenticated(approveIncorporationCertificate))
	http.Handle("/reject_incorporationCert/{organizationId}", isAuthenticated(rejectIncorporationCertificate))
	http.Handle("/approve_businessPremiseCert/{organizationId}", isAuthenticated(approveBusinessPremiseCertificate))
	http.Handle("/reject_businessPremiseCert/{organizationId}", isAuthenticated(rejectBusinessPremiseCertificate))
	http.Handle("/approve_companyLogo/{organizationId}", isAuthenticated(approveCompanyLogo))
	http.Handle("/reject_companyLogo/{organizationId}", isAuthenticated(rejectCompanyLogo))
	http.Handle("/approve_formc07/{organizationId}", isAuthenticated(approveFormC07))
	http.Handle("/reject_formc07/{organizationId}", isAuthenticated(rejectFormC07))
	http.Handle("/approve_idDoc/{organizationId}", isAuthenticated(approveIDDocumennt))
	http.Handle("/reject_idDoc/{organizationId}", isAuthenticated(rejectIDDocument))

	http.Handle("/update_memorandum", isAuthenticated(updateMemorandum))
	http.Handle("/update_memorandum_handler", isAuthenticated(updateMemorandumHandler))
	http.Handle("/update_coverLetter", isAuthenticated(updateCoverLetter))
	http.Handle("/update_coverLetter_handler", isAuthenticated(updateCoverLetterHandler))
	http.Handle("/update_businessCert", isAuthenticated(updateBusinessCertificate))
	http.Handle("/update_businessCert_handler", isAuthenticated(updateBusinessCertificateHandler))
	http.Handle("/update_incorporationCert", isAuthenticated(updateIncorporationCertificate))
	http.Handle("/update_incorporationCert_handler", isAuthenticated(updateIncorporationCertificateHandler))
	http.Handle("/update_businessPremiseCert", isAuthenticated(updateBusinessPremiseCertificate))
	http.Handle("/update_businessPremiseCert_handler", isAuthenticated(updateBusinessPremiseCertificateHandler))
	http.Handle("/update_passportPhoto", isAuthenticated(updateCompanyLogo))
	http.Handle("/update_passportPhoto_handler", isAuthenticated(updateCompanyLogoHandler))
	http.Handle("/update_formc07", isAuthenticated(updateFormC07))
	http.Handle("/update_formc07_handler", isAuthenticated(updateFormC07Handler))
	http.Handle("/update_idDoc", isAuthenticated(updateIDDocument))
	http.Handle("/update_idDoc_handler", isAuthenticated(updateIDDocumentHandler))
	http.Handle("/profile", isAuthenticated(profile))
	http.Handle("/profile_update", isAuthenticated(profileUpdate))
	http.Handle("/profile_update_handler", isAuthenticated(profileUpdateHandler))

	fmt.Println("Server Running on Port :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println(currentUser)
	if currentUser.Role == "user" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	} else if currentUser.Role == "admin" {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	} else if currentUser.Role == "superAdmin" {
		http.Redirect(w, r, "/approval_admin", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func profileUpdate(w http.ResponseWriter, r *http.Request) {
	if currentUser.Role != "user" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	type PageResult struct {
		User User
		QA   PasswordResetQuestion
	}
	var user User
	var securityQa PasswordResetQuestion
	userResp, err := findUserById(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	user = userResp.Body[0]
	secQaResp, err := findSecurityQAByUser(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(secQaResp.Body) == 0 {
		securityQa = PasswordResetQuestion{}
	} else {
		securityQa = secQaResp.Body[0]
	}
	p := PageResult{
		User: user,
		QA:   securityQa,
	}
	tmpl.ExecuteTemplate(w, "profile_update.html", p)
}

func profileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if currentUser.Role != "user" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	id := r.FormValue("userId")
	// email := strings.ToLower(r.FormValue("email"))
	question := r.FormValue("security_question")
	answer := strings.ToLower(r.FormValue("security_answer"))
	passwordResetQA := PasswordResetQuestion{}
	userResp, err := findUserById(dbFindUrl, currentUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	user := userResp.Body[0]

	passworedResetQaResp, err := findSecurityQAByUser(dbFindUrl, user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(passworedResetQaResp.Body) == 0 {
		passwordResetQA = PasswordResetQuestion{
			ID:       uuid.NewString(),
			Question: question,
			Answer:   answer,
			UserId:   id,
			Doctype:  "securityQa",
		}
		jsonData, err := json.Marshal(passwordResetQA)
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

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	} else {
		passwordResetQA = passworedResetQaResp.Body[0]
		passwordResetQA.Answer = answer
		passwordResetQA.Question = question

		jsonData, err := json.Marshal(&passwordResetQA)
		if err != nil {
			fmt.Println(err)
			return
		}
		request, err := http.NewRequest("PUT", dbUrl+"/"+passwordResetQA.ID, bytes.NewBuffer(jsonData))
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
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
