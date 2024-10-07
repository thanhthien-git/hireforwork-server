package interfaces

type ICompanyFilter struct {
	CompanyName   string `bson:"companyName" json:"companyName"`
	CompanyEmail  string `bson:"companyEmail" json:"companyEmail"`
	TypeOfCompany string `bson:"typeOfCompany" json:"typeOfCompany"`
}
