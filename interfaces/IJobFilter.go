package interfaces

type IJobFilter struct {
	JobTitle        string `json:"job_title,omitempty"`
	WorkingLocation string `json:"working_location,omitempty"`
	JobCategory     string `json:"job_category,omitempty"`
	CompanyName     string `json:"company_name,omitempty"`
}
