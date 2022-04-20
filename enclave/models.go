package enclave

type EnrolmentKeyState struct {
	Id           string   `tfdsk:"id"`
	Type         string   `tfdsk:"type"`
	ApprovalMode string   `tfdsk:"approval_mode"`
	Description  string   `tfdsk:"description"`
	Tags         []string `tfdsk:"tags"`
}
