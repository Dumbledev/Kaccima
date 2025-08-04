package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func admin(w http.ResponseWriter, r *http.Request) {
	type PageResult struct {
		UserCount          int
		PendingOrgCount    int
		PendingOrg         []Organization
		PendingBankPayment []BankTransfer
	}
	pendingBankTransfer := []BankTransfer{}
	pendingOrg := []Organization{}
	users := []User{}
	userResponse, err := findUsers(dbFindUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(userResponse.Body) != 0 {
		users = userResponse.Body
	}
	orgPendingStatusResponse, err := findOrganizationApprovalStatus(dbFindUrl, "Pending")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgPendingStatusResponse.Body) != 0 {
		pendingOrg = orgPendingStatusResponse.Body
	}

	bankTransferPendingStatusResponse, err := findBankPaymentApprovalStatus(dbFindUrl, "Pending")
	if err != nil {
		fmt.Println(err, "Err")
		return
	}
	if len(bankTransferPendingStatusResponse.Body) != 0 {
		pendingBankTransfer = bankTransferPendingStatusResponse.Body
	}

	p := PageResult{
		UserCount:          len(users),
		PendingOrgCount:    len(pendingOrg),
		PendingOrg:         pendingOrg,
		PendingBankPayment: pendingBankTransfer,
	}
	tmpl.ExecuteTemplate(w, "admin.html", p)
}

func adminDocuments(w http.ResponseWriter, r *http.Request) {
	organizations := []Organization{}
	orgResp, err := findAllOrganizations(dbFindUrl)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organizations = orgResp.Body
	}
	tmpl.ExecuteTemplate(w, "admin_documents.html", organizations)
}

func adminOrganizationDocuments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, id)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	tmpl.ExecuteTemplate(w, "admin_organization_documents.html", organization)
}

func adminMembers(w http.ResponseWriter, r *http.Request) {
	organizations := []Organization{}
	accceptedOrganizations := []Organization{}
	orgResp, err := findAllOrganizations(dbFindUrl)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organizations = orgResp.Body
	}
	for _, v := range organizations {
		if v.Status == "Pending" {
			continue
		} else {
			accceptedOrganizations = append(accceptedOrganizations, v)
		}
	}
	tmpl.ExecuteTemplate(w, "admin_members.html", accceptedOrganizations)
}

func adminPayment(w http.ResponseWriter, r *http.Request) {
	var bankTransfers []BankTransfer
	bankTransferPendingStatusResponse, err := findOrganizationPayments(dbFindUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(bankTransferPendingStatusResponse.Body) == 0 {
		bankTransfers = []BankTransfer{}
		fmt.Println("No Record Found")
	} else {
		bankTransfers = bankTransferPendingStatusResponse.Body
	}
	tmpl.ExecuteTemplate(w, "admin_payments.html", bankTransfers)
}

func adminReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("organizationId")
	orgResp, err := findOrganizationById(dbFindUrl, id)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) == 0 {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}
	organization := orgResp.Body[0]
	tmpl.ExecuteTemplate(w, "admin_report.html", organization)
}

func adminReportHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("organizationId")
	reportType := r.FormValue("reportType")
	reason := r.FormValue("reportReason")
	date := r.FormValue("reportDate")
	fmt.Println(id, "id")
	orgResp, err := findOrganizationById(dbFindUrl, id)
	fmt.Println(orgResp.Body)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) == 0 {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}

	org := orgResp.Body[0]
	report := Report{
		ID:               uuid.NewString(),
		Type:             reportType,
		Reason:           reason,
		Date:             date,
		OrganizationId:   id,
		OrganizationName: org.Name,
		Doctype:          "report",
	}

	jsonData, err := json.Marshal(report)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(jUser))
	request, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Server Error(1)")
		return
	}
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Server Error(2)")
		return
	}
	defer res.Body.Close()

	http.Redirect(w, r, "/admin_members", http.StatusPermanentRedirect)
}

func adminSettings(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin_settings.html", nil)
}

func acceptReceipt(w http.ResponseWriter, r *http.Request) {
	paymentId := r.PathValue("paymentId")
	fmt.Println(paymentId, "id")
	payment := BankTransfer{}
	paymentResp, err := findOrganizationPaymentByID(dbFindUrl, paymentId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(paymentResp.Body) != 0 {
		payment = paymentResp.Body[0]
	}
	payment.Status = "Accepted"
	fmt.Println(payment)

	jsonData, err := json.Marshal(&payment)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+payment.ID, bytes.NewBuffer(jsonData))
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
	http.Redirect(w, r, "/approval_admin", http.StatusPermanentRedirect)
}

func rejectReceipt(w http.ResponseWriter, r *http.Request) {
	paymentId := r.PathValue("paymentId")
	payment := BankTransfer{}
	paymentResp, err := findOrganizationPaymentByID(dbFindUrl, paymentId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(paymentResp.Body) != 0 {
		payment = paymentResp.Body[0]
	}
	payment.Status = "Rejected"

	jsonData, err := json.Marshal(&payment)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+payment.ID, bytes.NewBuffer(jsonData))
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
	http.Redirect(w, r, "/approval_admin", http.StatusPermanentRedirect)
}

func acceptOrganization(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	org := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		org = orgResp.Body[0]
	}
	org.Status = "Accepted"

	jsonData, err := json.Marshal(&org)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+org.ID, bytes.NewBuffer(jsonData))
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
	http.Redirect(w, r, "/approval_admin", http.StatusPermanentRedirect)
}

func rejectOrganization(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	org := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		org = orgResp.Body[0]
	}
	org.Status = "Rejected"

	jsonData, err := json.Marshal(&org)
	if err != nil {
		fmt.Println(err)
		return
	}
	request, err := http.NewRequest("PUT", dbUrl+"/"+org.ID, bytes.NewBuffer(jsonData))
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
	http.Redirect(w, r, "/approval_admin", http.StatusPermanentRedirect)
}

func approveMemorandum(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.MemorandumApproval = "Approved"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func rejectMemorandum(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.MemorandumApproval = "Rejected"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func approveCoverLetter(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.CoverLetterApproval = "Approved"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func rejectCoverLetter(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.CoverLetterApproval = "Rejected"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func approveBusinessCertificate(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.BusinessCertificateApproval = "Approved"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func rejectBusinessCertificate(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.BusinessCertificateApproval = "Rejected"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func approveIncorporationCertificate(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.IncorporationCertificateApproval = "Approved"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func rejectIncorporationCertificate(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.IncorporationCertificateApproval = "Rejected"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func approveBusinessPremiseCertificate(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.BusinessPremiseCertificateApproval = "Approved"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func rejectBusinessPremiseCertificate(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.BusinessPremiseCertificateApproval = "Rejected"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func approveCompanyLogo(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.CompanyLogo = "Approved"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func rejectCompanyLogo(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.CompanyLogo = "Rejected"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func approveFormC07(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.FormC07Approval = "Approved"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func rejectFormC07(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.FormC07Approval = "Rejected"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func approveIDDocumennt(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.IDDocumentApproval = "Approved"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func rejectIDDocument(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	organization := Organization{}
	orgResp, err := findOrganizationById(dbFindUrl, organizationId)
	if err != nil {
		tmpl.ExecuteTemplate(w, "500.html", "Error")
		return
	}
	if len(orgResp.Body) != 0 {
		organization = orgResp.Body[0]
	}
	organization.IDDocumentApproval = "Rejected"

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
	http.Redirect(w, r, "/admin_organization_documents/"+organizationId, http.StatusPermanentRedirect)
}

func approvalAdmin(w http.ResponseWriter, r *http.Request) {
	type PageResult struct {
		UserCount          int
		PendingOrgCount    int
		PendingOrg         []Organization
		PendingBankPayment []BankTransfer
	}
	pendingBankTransfer := []BankTransfer{}
	pendingOrg := []Organization{}
	users := []User{}
	userResponse, err := findUsers(dbFindUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(userResponse.Body) != 0 {
		users = userResponse.Body
	}
	orgPendingStatusResponse, err := findOrganizationApprovalStatus(dbFindUrl, "Pending")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(orgPendingStatusResponse.Body) != 0 {
		pendingOrg = orgPendingStatusResponse.Body
	}

	bankTransferPendingStatusResponse, err := findBankPaymentApprovalStatus(dbFindUrl, "Pending")
	if err != nil {
		fmt.Println(err, "Err")
		return
	}
	if len(bankTransferPendingStatusResponse.Body) != 0 {
		pendingBankTransfer = bankTransferPendingStatusResponse.Body
	}

	p := PageResult{
		UserCount:          len(users),
		PendingOrgCount:    len(pendingOrg),
		PendingOrg:         pendingOrg,
		PendingBankPayment: pendingBankTransfer,
	}
	tmpl.ExecuteTemplate(w, "approval_admin.html", p)
}
