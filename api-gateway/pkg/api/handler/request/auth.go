package request

type SignUp struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"min=6,max=25"`
}

type SignIn struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"min=6,max=25"`
}
