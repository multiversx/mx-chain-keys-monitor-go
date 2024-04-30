package core

// ValidatorStatistics represents the DTO returned by the API
type ValidatorStatistics struct {
	TempRating float32 `json:"tempRating"`
	Rating     float32 `json:"rating"`
}

// ValidatorStatisticsResponse represents the DTO for the validator/statistics response
type ValidatorStatisticsResponse struct {
	Data  map[string]map[string]*ValidatorStatistics `json:"data"`
	Error string                                     `json:"error"`
	Code  string                                     `json:"code"`
}

// CheckResponse defines the checking response DTO
type CheckResponse struct {
	HexBLSKey string
	Status    string
}

// OutputMessage defines the message to be sent to an output notifier
type OutputMessage struct {
	Type               MessageOutputType
	IdentifierType     string
	Identifier         string
	ShortIdentifier    string
	IdentifierURL      string
	ExecutorName       string
	ProblemEncountered string
}

// Address defines the address DTO with it's 2 representations
type Address struct {
	Hex    string
	Bech32 string
}
