package domain

type SourceCount struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

type FacultyCount struct {
	Faculty string `json:"faculty"`
	Count   int    `json:"count"`
}

type AgeCount struct {
	Age   int `json:"age"`
	Count int `json:"count"`
}

type FacultyRegisterCount struct {
	Faculty string
	Count   int
}

type FacultyPercent struct {
	Faculty        string  `json:"faculty"`
	FirstInterest  float64 `json:"first_interest"`
	SecondInterest float64 `json:"second_interest"`
	ThirdInterest  float64 `json:"third_interest"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type StudentProfile struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	FirstInterest  string `json:"first_interest"`
	SecondInterest string `json:"second_interest"`
	ThirdInterest  string `json:"third_interest"`
}

type AttendedCount struct {
	Count int `json:"count"`
}
