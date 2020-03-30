package dto

// Answered statistics
type Answered struct {
	Lost  int `json:"lost"`
	Won   int `json:"won"`
	Total int `json:"total"`
}

// Wrapper to data
type Wrapper struct {
	Data interface{} `json:"data"`
}
