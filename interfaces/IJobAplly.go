package interfaces

type IJobApply struct {
	JobID       string `json:"jobID"`
	IDCareer    string `json:"careerID"`
	CompanyID   string `json:"companyID"`
	CareerCV    string `json:"careerCV"`
	CareerEmail string `json:"careerEmail"`
}
