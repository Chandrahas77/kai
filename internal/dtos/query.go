package dtos

type FliterRequest struct {
	Filters struct {
		Severity string `json:"severity"`
	} `json:"filters"`
}
