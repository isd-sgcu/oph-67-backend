package domain

type SourceCount struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

type FacultyCount struct {
	Faculty string `json:"faculty"`
	Count   int    `json:"count"`
}

type FacultyPercent struct {
	Faculty        string  `json:"faculty"`
	FirstInterest  float64 `json:"first_interest"`
	SecondInterest float64 `json:"second_interest"`
	ThirdInterest  float64 `json:"third_interest"`
}
