package models

type Performance struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Artists   []Artist `json:"artists"`
	Stage     Stage    `json:"stage"`
	Date      string   `json:"date"`
	Day       string   `json:"day"`
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
}
