package music

type Genre struct {
	Id          int64  `json:"-"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
