package interfaces

type IJobFilter struct {
	JobTitle        string   `json:"jobTitle"`
	CompanyName     string   `json:"companyName"`
	DateCreateFrom  string   `json:"dateCreateFrom"`
	DateCreateTo    string   `json:"dateCreateTo"`
	EndDateFrom     string   `json:"endDateFrom"`
	EndDateTo       string   `json:"endDateTo"`
	SalaryFrom      int64    `json:"salaryFrom"`
	SalaryTo        int64    `json:"salaryTo"`
	WorkingLocation []string `json:"workingLocation"`
	JobRequirement  []string `json:"jobRequirement"`
	JobCategory     []string `json:"jobCategory"`
	JobLevel        string   `json:"jobLevel"`
	IsHot           bool     `json:"isHot"`
	Query           string   `json:"query"`
}
