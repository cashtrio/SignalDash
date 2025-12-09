package dto

type OkResponse struct {
	Message string `json:"message"`
}

type ScrapeDollarIndexResponse struct {
	DollarIndex string `json:"dollar_index"`
}
