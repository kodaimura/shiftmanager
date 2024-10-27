package request


type PostShiftPreferred struct {
	Dates *string `json:"dates"`
	Notes *string `json:"notes"`
}

type ShiftPreferredUri struct {
    Year  int `uri:"year" binding:"required"`
    Month int `uri:"month" binding:"required"`
}

type ShiftPreferredQuery struct {
    Year  int `form:"year" binding:"required"`
    Month int `form:"month" binding:"required"`
}