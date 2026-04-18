package api

type HourlyCheckinSummary struct {
	Hour  int `json:"hour"`
	Count int `json:"count"`
}

type LinkCount struct {
	Link  string `json:"link"`
	Count int    `json:"count"`
}
