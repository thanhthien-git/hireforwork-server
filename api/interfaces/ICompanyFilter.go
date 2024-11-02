package interfaces

type ICompanyFilter struct {
	CompanyName  string `bson:"companyName" json:"companyName"`
	CompanyEmail string `bson:"companyEmail" json:"companyEmail"`
	StartDate    *string
	EndDate      *string
}
