package user_request

type SaveOne struct {
	Id    string `json:"id" binding:"max=128"`
	Name  string `json:"name" binding:"required,max=255"`
	Email string `json:"email" binding:"required,email,max=255"`
}
