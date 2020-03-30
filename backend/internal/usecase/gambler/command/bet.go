package command

// PlaceBet command
type PlaceBet struct {
	Topic      string `json:"topic"`      // ID
	User       string `json:"-"`          // UserID fetched via BE
	Prediction string `json:"prediction"` // Bet
	Stake      int    `json:"stake"`      // Reputation
}
