package main

import (
	"fmt"
	"log"
	"net/http"

	"html/template"

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

	http.HandleFunc("/", index)
	http.HandleFunc("/kaccimanew2025", adminSignUp)
	http.HandleFunc("/kaccimanew2025_handler", adminSignUpHandler)

	http.HandleFunc("/sign_up", signUp)
	http.HandleFunc("/sign_up_handler", signUpHandler)
	http.HandleFunc("/sign_in", signIn)
	http.HandleFunc("/sign_in_handler", signInHandler)
	http.HandleFunc("/sign_out", signOut)
	http.HandleFunc("/forgot_password", forgotPassword)

	http.Handle("/dashboard", isAuthenticated(organizationDashboard))
	http.Handle("/notification", isAuthenticated(notification))
	// http.Handle("/profile", isAuthenticated(profile))
	http.Handle("/payment", isAuthenticated(payment))
	http.Handle("/bank_transfer_handler", isAuthenticated(bankTransferHandler))
	http.Handle("/reviewedDocuments", isAuthenticated(reviewedDocuments))
	http.Handle("/organization", isAuthenticated(organization))
	http.Handle("/organization_register", isAuthenticated(organizationRegister))
	http.Handle("/organization_register_handler", isAuthenticated(organizationRegisterHandler))

	http.Handle("/admin", isAuthenticated(admin))
	http.Handle("/admin_documents", isAuthenticated(adminDocuments))
	http.Handle("/admin_organization_documents/{organizationId}", isAuthenticated(adminOrganizationDocuments))
	http.Handle("/admin_members", isAuthenticated(adminMembers))
	http.Handle("/admin_payments", isAuthenticated(adminPayment))
	http.Handle("/admin_report", isAuthenticated(adminReport))
	http.Handle("/admin_settings", isAuthenticated(adminSettings))

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

	http.Handle("/update_rejected_documents", isAuthenticated(updateRejectedDocuments))
	// http.Handle("/update_rejected_documents_handler", isAuthenticated(updateRejectedDocumentsHandler))
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
	http.Handle("/update_passportPhoto", isAuthenticated(updatePassportPhoto))
	http.Handle("/update_passportPhoto_handler", isAuthenticated(updatePassportPhotoHandler))
	http.Handle("/update_formc07", isAuthenticated(updateFormC07))
	http.Handle("/update_formc07_handler", isAuthenticated(updateFormC07Handler))
	http.Handle("/update_idDoc", isAuthenticated(updateIDDocument))
	http.Handle("/update_idDoc_handler", isAuthenticated(updateIDDocumentHandler))
	// http.Handle("/profile", isAuthenticated(profile))
	// http.Handle("/profile_update", isAuthenticated(profileUpdate))
	// http.Handle("/profile_update_handler", isAuthenticated(profileUpdateHandler))

	fmt.Println("Server Running on Port :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

// func profile(w http.ResponseWriter, r *http.Request) {
// 	var profile UserProfile
// 	profileResp, err := findUserProfile(dbFindUrl, currentUser.ID)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	profile = profileResp.Body[0]
// 	tmpl.ExecuteTemplate(w, "profile.html", profile)
// }

// func profileUpdate(w http.ResponseWriter, r *http.Request) {
// 	var profile UserProfile
// 	profileResp, err := findUserProfile(dbFindUrl, currentUser.ID)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	profile = profileResp.Body[0]
// 	tmpl.ExecuteTemplate(w, "profile_update.html", profile)
// }

// func profileUpdateHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	// profileId := r.FormValue("profileId")
// 	http.Redirect(w, r, "/profile", http.StatusPermanentRedirect)
// }
