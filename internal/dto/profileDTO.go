package dto

type CreateProfileReq struct {
	Name *string `json:"name"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ProcessedData struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	SampleSize  int     `json:"sample_size"`
	IsConfident bool    `json:"is_confident"`
	ProcessedAt string  `json:"processed_at"`
}

type SuccessResponse struct {
	Status string        `json:"status"`
	Data   ProcessedData `json:"data"`
}

type CountriesResponse struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type GenderizeRawData struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

type AgifyRawData struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Count int    `json:"count"`
}

type NationalizeRawData struct {
	Name    string              `json:"name"`
	Count   int                 `json:"count"`
	Country []CountriesResponse `json:"country"`
}
