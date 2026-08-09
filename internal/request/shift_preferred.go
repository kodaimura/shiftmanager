package request

type PostShiftPreferred struct {
	Dates *string `json:"dates" binding:"omitempty,max=100"`
	Notes *string `json:"notes" binding:"omitempty,max=500"`
}

type ShiftPreferredUri struct {
	Year  int `uri:"year" binding:"required,gte=2000,lte=2100"`
	Month int `uri:"month" binding:"required,gte=1,lte=12"`
}

type ShiftPreferredQuery struct {
	Year  int `form:"year" binding:"required,gte=2000,lte=2100"`
	Month int `form:"month" binding:"required,gte=1,lte=12"`
}
