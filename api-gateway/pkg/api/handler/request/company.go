package request

type CompanyRequest struct {
	Name           string `json:"name" binding:"required,min=3,max=15"`
	CEO            string `json:"ceo" binding:"required,min=3,max=25"`
	TotalEmployees int    `json:"total_employees" binding:"required,min=30"`
}
