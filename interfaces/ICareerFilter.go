package interfaces

type ICareerFilter struct {
	FirstName string `json: "careerFirstName"`
	LastName  string `json: "lastName"`
	Email     string `json:"email"`
	Phone     string `json:"careerPhone"`
}
