package request


type PostShiftPreferred struct {
	Dates *string `json:"dates"`
	Notes *string `json:"notes"`
}

type GetShiftPreferredParams struct {
    Year  int `uri:"year" binding:"required"`
    Month int `uri:"month" binding:"required"`
}

type PostShiftPreferredParams struct {
    Year  int `uri:"year" binding:"required"`
    Month int `uri:"month" binding:"required"`
}