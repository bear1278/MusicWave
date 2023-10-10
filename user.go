package music

type User struct {
	Id       int64  `json:"-"`
	Name     string `json:"name" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
