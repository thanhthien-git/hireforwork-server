package interfaces

type IJobFilter struct {
	JobTitle        string   `json:"jobTitle"`
	CompanyName     string   `json:"companyName"`
	DateCreateFrom  string   `json:"dateCreateFrom"`
	DateCreateTo    string   `json:"dateCreateTo"`
	EndDateFrom     string   `json:"endDateFrom"`
	EndDateTo       string   `json:"endDateTo"`
	SalaryFrom      string   `json:"salaryFrom"`
	SalaryTo        string   `json:"salaryTo"`
	WorkingLocation []string `json:"workingLocation"`
	JobRequirement  []string `json:"jobRequirement"`
	JobLevel        string   `json:"jobLevel"`
}
