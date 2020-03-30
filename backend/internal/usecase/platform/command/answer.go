package command

import "strings"

// Answer a topic
type Answer struct {
	Topic      string   `json:"topic"`
	Answer     string   `json:"answer"`
	Variations []string `json:"variations"`
}

// IsTrue check prediction against answer
func (ans *Answer) IsTrue(prediction string) bool {
	prediction = strings.TrimSpace(prediction)
	answer := strings.TrimSpace(ans.Answer)
	if prediction == answer {
		return true
	}

	for _, variety := range ans.Variations {
		variety = strings.TrimSpace(variety)
		if prediction == variety {
			return true
		}
	}

	return false
}
