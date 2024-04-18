package core

// ValidatorStatistics represents the DTO returned by the API
type ValidatorStatistics struct {
	TempRating float32 `json:"tempRating"`
	Rating     float32 `json:"rating"`
}

// CheckResponse defines the checking response DTO
type CheckResponse struct {
	HexBLSKey string
	Status    string
}
