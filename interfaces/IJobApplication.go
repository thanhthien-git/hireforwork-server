package interfaces

type IJobApplicationFilter struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	CareerEmail string `json:"careerEmail"`
	CreateFrom  string `json:"createFrom"`
	CreateTo    string `json:"createTo"`
	JobLevel    string `json:"jobLevel"`
	JobTitle    string `json:"jobTitle"`
	Status      string `json:"status"`
}
